import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
from pathlib import Path

forDb = False
files = {
    20:    Path("Long/Full/Id/getById_checkpoint_results_20.xlsx"),
    2000:  Path("Long/Full/Id/getById_checkpoint_results_2000.xlsx"),
    8000:  Path("Long/Full/Id/getById_checkpoint_results_8000.xlsx"),
    16000: Path("Long/Full/Id/getById_checkpoint_results_16000.xlsx"),
}
sheet_name = "Results"
groups = ["First", "Center", "Last"]

DB_COLS = {
    "First":  "First - Db duration",
    "Center": "Center - Db duration",
    "Last":   "Last - Db duration",
}
BC_COLS = {
    "First":  "First - Blockchain duration",
    "Center": "Center - Blockchain duration",
    "Last":   "Last - Blockchain duration",
}

out_dir = Path("Long_plots/Full/id")
out_dir.mkdir(parents=True, exist_ok=True)

def parse_duration(val):
    """'1.23ms' -> 1.23 ; '1.5s' -> 1500.0 ; liczby -> ms"""
    if isinstance(val, str):
        v = val.strip().lower().replace(" ", "").replace(",", ".")
        if v.endswith("ms"):
            return float(v[:-2])
        if v.endswith("s"):
            return float(v[:-1]) * 1000.0
        return float(v)
    return float(val)

def extract_data(xlsx, colmap):
    """Wczytaj jeden plik i zwróć listę [First, Center, Last] z wartościami w ms"""
    df = pd.read_excel(xlsx, sheet_name=sheet_name)
    df = df.loc[:, ~df.columns.astype(str).str.startswith("Unnamed")]  # usuń puste kolumny
    values = []
    for g in groups:
        col = colmap[g]
        vals = df[col].apply(parse_duration).values
        values.append(vals)
    return values

def plot_single(data_list, ylabel, out_png, blocks, which):
    """Rysuje boxplot z trzema grupami (First/Center/Last) + jitter + średnia + std dev
       i wypisuje statystyki do konsoli."""
    plt.figure(figsize=(6, 6))
    plt.boxplot(
        data_list,
        patch_artist=True,
        labels=groups,
        showfliers=False,
        widths=0.5
    )

    print(f"\n📊 Statystyki dla {which}, {blocks} bloków:")
    for i, (y, g) in enumerate(zip(data_list, groups), start=1):
        x = np.random.normal(i, 0.08, size=len(y))  # jitter ±0.08
        plt.scatter(x, y, alpha=0.6, s=20)

        mean = np.mean(y)
        std  = np.std(y)
        median = np.median(y)
        q1 = np.percentile(y, 25)
        q3 = np.percentile(y, 75)
        min_val = np.min(y)
        max_val = np.max(y)

        plt.plot(i, mean, 'ro', markersize=7, label="Mean" if i == 1 else "")
        plt.errorbar(i, mean, yerr=std, fmt='none', ecolor='r', capsize=5,
                     label="Std Dev" if i == 1 else "")

        # wypisz statystyki
        print(f"  {g}:")
        print(f"    Min     = {min_val:.3f} ms")
        print(f"    Q1      = {q1:.3f} ms")
        print(f"    Median  = {median:.3f} ms")
        print(f"    Q3      = {q3:.3f} ms")
        print(f"    Max     = {max_val:.3f} ms")
        print(f"    Mean    = {mean:.3f} ms")
        print(f"    Std Dev = {std:.3f} ms")

    plt.ylabel(ylabel)
    # plt.xlabel("Document positions in the database")
    plt.xlabel("Position of the document hash in the blockchain")
    plt.legend()
    plt.tight_layout()
    plt.savefig(out_png, dpi=300, bbox_inches="tight")
    plt.close()

for blocks, path in files.items():
    mongo_data = extract_data(path, DB_COLS)
    eth_data   = extract_data(path, BC_COLS)

    if forDb:
        plot_single(
            mongo_data,
            "Duration [ms]",
            out_dir / f"mongodb_{blocks}.png",
            blocks,
            "MongoDB"
        )
    plot_single(
        eth_data,
        "Duration [ms]",
        out_dir / f"ethereum_{blocks}.png",
        blocks,
        "Ethereum"
    )

print("\n✅ Zapisano wykresy w folderze:", out_dir)

import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
from pathlib import Path

files = {
    10: Path("Time\send_results_10minutes.xlsx"),
    5: Path("Time\send_results_5minutes.xlsx"),
    1: Path("Time\send_results_1minute.xlsx"),
    15: Path("Time\send_results_15seconds.xlsx"),
}
sheet_name = "Results"

out_dir = Path("Time_plots\send")
out_dir.mkdir(parents=True, exist_ok=True)

def parse_duration(val):
    """Konwertuje '123.4ms' -> 123.4 ; '1.2s' -> 1200.0 ; liczby -> ms"""
    if isinstance(val, str):
        v = val.strip().lower().replace(",", ".")
        if v.endswith("ms"):
            return float(v.replace("ms", ""))
        elif v.endswith("s"):
            return float(v.replace("s", "")) * 1000
    return float(val)

def make_scatter(xlsx_path: Path, blocks: int):
    df = pd.read_excel(xlsx_path, sheet_name=sheet_name)
    df["Blockchain_ms"] = df["Blockchain duration"].apply(parse_duration)
    df["Tx_index"] = np.arange(1, len(df) + 1)

    mean_val = df["Blockchain_ms"].mean()
    median_val = df["Blockchain_ms"].median()

    plt.figure(figsize=(10,6))
    plt.scatter(df["Tx_index"], df["Blockchain_ms"], alpha=0.6, s=12, color="orange")
    plt.axhline(mean_val, color="red", linestyle="--", label=f"Mean = {mean_val:.2f} ms")
    plt.axhline(median_val, color="blue", linestyle="--", label=f"Median = {median_val:.2f} ms")

    plt.xlabel("Transaction index")
    plt.ylabel("Duration [ms]")
    plt.legend()
    plt.tight_layout()

    out_path = out_dir / f"send_{blocks}.png"
    plt.savefig(out_path, dpi=400, bbox_inches="tight")
    plt.close()
    print(f"✅ Saved: {out_path}")

for blocks, path in files.items():
    make_scatter(path, blocks)

print("\nAll scatter plots generated.")

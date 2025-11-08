import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
from pathlib import Path
import re

# =========================
# CONFIG
# =========================
FILES = [
    Path("Time_ef/Light/audit_results_Light_1200_15s.xlsx"),
    Path("Time_ef/Light/audit_results_Light_1200_1min.xlsx"),
    Path("Time_ef/Light/audit_results_Light_1200_5min.xlsx"),
    Path("Time_ef/Light/audit_results_Light_1200_10min.xlsx"),
]
SHEET_NAME = "Results"
COL_BC = "Blockchain duration"

OUTPUT_DIR = Path("Time_ef_plot")
OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

# =========================
# HELPERS
# =========================
def pick_sheet(xlsx_path: Path, preferred: str | None) -> str:
    xls = pd.ExcelFile(xlsx_path)
    if preferred and preferred in xls.sheet_names:
        return preferred
    return xls.sheet_names[0]

def parse_duration(val):
    """'123.4ms' -> 123.4 ; '1.2s' -> 1200.0 ; numbers -> ms"""
    if isinstance(val, str):
        v = val.strip().lower().replace(",", ".")
        if v.endswith("ms"):
            return float(v[:-2])
        if v.endswith("s"):
            return float(v[:-1]) * 1000.0
        return float(v)
    return float(val)

def load_ms_from_file(xlsx_path: Path, col_name: str) -> np.ndarray:
    sheet = pick_sheet(xlsx_path, SHEET_NAME)
    df = pd.read_excel(xlsx_path, sheet_name=sheet)
    if col_name not in df.columns:
        raise ValueError(f"{xlsx_path.name}: missing column '{col_name}'. Available: {list(df.columns)}")
    return df[col_name].apply(parse_duration).to_numpy()

def extract_block_time_label(path: Path) -> str:
    """
    Map filename suffix to pretty label:
      ..._15s.xlsx   -> '15 s'
      ..._1min.xlsx  -> '1 m'
      ..._5min.xlsx  -> '5 m'
      ..._10min.xlsx -> '10 m'
    Fallback: last number from stem.
    """
    s = path.stem.lower()
    if s.endswith("15s"):
        return "15s"
    if s.endswith("1min"):
        return "1min"
    if s.endswith("5min"):
        return "5min"
    if s.endswith("10min"):
        return "10min"
    nums = re.findall(r"\d+", path.stem)
    return nums[-1] if nums else path.stem

def make_multi_boxplot(values_list: list[np.ndarray],
                       x_tick_labels: list[str],
                       ylabel: str,
                       x_label: str,
                       out_path: Path,
                       title: str | None = None,
                       scatter_alpha: float = 0.6,
                       jitter: float = 0.06,
                       log_scale: bool = False):
    """
    Draw one plot with multiple boxplots (each from one file).
    Adds jittered points, mean (dot), and standard deviation (error bar).
    Style matches the provided reference snippet.
    """
    plt.figure(figsize=(7.5, 6))
    plt.boxplot(values_list, patch_artist=True, labels=x_tick_labels,
                showfliers=False, widths=0.6)

    # points + stats
    print(f"\n📊 Stats for: {out_path.name}")
    for i, (vals, lab) in enumerate(zip(values_list, x_tick_labels), start=1):
        x = np.random.normal(i, jitter, size=len(vals))
        plt.scatter(x, vals, alpha=scatter_alpha, s=18)

        mean = float(np.mean(vals))
        std  = float(np.std(vals))
        med  = float(np.median(vals))
        q1   = float(np.percentile(vals, 25))
        q3   = float(np.percentile(vals, 75))
        vmin = float(np.min(vals))
        vmax = float(np.max(vals))

        plt.plot(i, mean, 'ro', markersize=7, label="Mean" if i == 1 else "")
        plt.errorbar(i, mean, yerr=std, fmt='none', ecolor='r', capsize=5,
                     label="Std Dev" if i == 1 else "")

        print(f"  Block time {lab}:")
        print(f"    Min     = {vmin:.3f} ms")
        print(f"    Q1      = {q1:.3f} ms")
        print(f"    Median  = {med:.3f} ms")
        print(f"    Q3      = {q3:.3f} ms")
        print(f"    Max     = {vmax:.3f} ms")
        print(f"    Mean    = {mean:.3f} ms")
        print(f"    Std Dev = {std:.3f} ms")

    if log_scale:
        plt.yscale("log")

    if title:
        plt.title(title)  # <- you asked for no title, so keep title=None in call
    plt.ylabel(ylabel)
    plt.xlabel(x_label)
    plt.legend()
    plt.tight_layout()
    plt.savefig(out_path.with_suffix(".png"), dpi=400, bbox_inches="tight")
    plt.close()
    print(f"[OK] Saved: {out_path.with_suffix('.png').name}")

# =========================
# MAIN
# =========================
xlabels = [extract_block_time_label(p) for p in FILES]

# Ethereum/Blockchain duration plot (4 boxplots: 15 s, 1 m, 5 m, 10 m)
eth_values = [load_ms_from_file(p, COL_BC) for p in FILES]
make_multi_boxplot(
    eth_values,
    x_tick_labels=xlabels,
    ylabel="Duration [ms]",
    x_label="Block time",
    out_path=OUTPUT_DIR / "ethereum_4boxplots",
    title=None  # no title
)

print("\n✅ Done. Plot saved in:", OUTPUT_DIR.as_posix())

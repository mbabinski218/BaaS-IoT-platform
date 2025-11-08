import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from pathlib import Path

# =========================
# CONFIG
# =========================
FILES = {
    "15 s": Path("Time_ef/Full/audit_results_Full_1200_15s.xlsx"),
    "1 m":  Path("Time_ef/Full/audit_results_Full_1200_1min.xlsx"),
    "5 m":  Path("Time_ef/Full/audit_results_Full_1200_5min.xlsx"),
    "10 m": Path("Time_ef/Full/audit_results_Full_1200_10min.xlsx"),
}
SHEET = "Results"

# output dir & file
OUTPUT_DIR = Path("Time_ef_plot")
OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
OUTPUT_PNG = OUTPUT_DIR / "boxplot_blockchain_duration_by_block_time.png"

# =========================
# HELPERS
# =========================
def parse_duration_ms(val):
    """'123.4ms' -> 123.4 ; '1.2s' -> 1200.0 ; numbers -> ms"""
    if pd.isna(val):
        return np.nan
    if isinstance(val, (int, float, np.integer, np.floating)):
        return float(val)
    v = str(val).strip().lower().replace(" ", "").replace(",", ".")
    try:
        if v.endswith("ms"):
            return float(v[:-2])
        if v.endswith("s"):
            return float(v[:-1]) * 1000.0
        return float(v)
    except ValueError:
        return np.nan

def find_blockchain_duration_column(df: pd.DataFrame) -> str:
    exact = [c for c in df.columns if str(c).strip().lower() == "blockchain duration"]
    if exact:
        return exact[0]
    candidates = [c for c in df.columns if "blockchain" in str(c).lower() and "duration" in str(c).lower()]
    if candidates:
        return candidates[0]
    raise KeyError(f"Missing 'Blockchain duration' column. Available: {list(df.columns)}")

# =========================
# LOAD & PREPARE DATA
# =========================
series_list = []
labels = []

for label, path in FILES.items():
    df = pd.read_excel(path, sheet_name=SHEET)
    df = df.loc[:, ~df.columns.astype(str).str.startswith("Unnamed")]
    col = find_blockchain_duration_column(df)
    vals_ms = df[col].apply(parse_duration_ms).dropna().astype(float).values
    series_list.append(vals_ms)
    labels.append(label)

# =========================
# PLOT (style like your grouped example: showfliers=False, mean marker + std bar)
# =========================
plt.figure(figsize=(10, 6))
base_pos = np.arange(len(series_list))  # 0..3

# single series, so no offsets; keep visual similar (widths, patch_artist, etc.)
bp = plt.boxplot(
    series_list,
    positions=base_pos,
    widths=0.5,
    patch_artist=True,
    showfliers=False,
    labels=[None]*len(base_pos)
)

# add mean and std per box (marker + errorbar)
for x, vals in zip(base_pos, series_list):
    if len(vals) == 0:
        continue
    mean = float(np.mean(vals))
    std  = float(np.std(vals))
    plt.plot(x, mean, 'o', markersize=5)
    plt.errorbar(x, mean, yerr=std, fmt='none', capsize=4)

# axes labels (English), no title
plt.xticks(base_pos, labels)
plt.xlabel("Block time")
plt.ylabel("Blockchain duration [ms]")

plt.tight_layout()
plt.savefig(OUTPUT_PNG, dpi=300, bbox_inches="tight")
plt.close()

print(f"Saved: {OUTPUT_PNG}")

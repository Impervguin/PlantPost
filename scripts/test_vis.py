import json
import os
import pandas as pd
import plotly.graph_objects as go
from plotly.subplots import make_subplots

# === Настройки ===
base_dir = "test_results"
subdirs = ["debug_no_tracing", "debug_with_tracing", "info_no_tracing", "info_with_tracing"]
output_html = "test_results_dashboard.html"

records = []

# === Считываем все JSON ===
for subdir in subdirs:
    path = os.path.join(base_dir, subdir)
    for filename in os.listdir(path):
        if filename.endswith(".json") and "summary_statistics" not in filename:
            with open(os.path.join(path, filename), "r") as f:
                data = json.load(f)
            rec = {
                "category": subdir,
                "iteration": data['iteration'],
                "avg_time": data["test_results"]["avg_time"],
                "min_time": data["test_results"]["min_time"],
                "max_time": data["test_results"]["max_time"],
                "cpu_avg": data["metrics"]["cpu_avg_value"],
                "cpu_min": data["metrics"]["cpu_min_value"],
                "cpu_max": data["metrics"]["cpu_max_value"],
                "memory_avg": data["metrics"]["memory_avg_value"] / (1024**2),
                "memory_min": data["metrics"]["memory_min_value"] / (1024**2),
                "memory_max": data["metrics"]["memory_max_value"] / (1024**2),
                "execution_times": data["test_results"]["execution_times"],
            }
            records.append(rec)

# === Превращаем в DataFrame ===
df = pd.DataFrame(records)

# === Средние значения по категориям ===
agg_df = df.groupby("category").agg({
    "avg_time": "mean",
    "min_time": "mean",
    "max_time": "mean",
    "cpu_avg": "mean",
    "cpu_min": "mean",
    "cpu_max": "mean",
    "memory_avg": "mean",
    "memory_min": "mean",
    "memory_max": "mean"
}).reset_index()

# === Форматируем данные для сводной таблицы до 2 знаков после запятой ===
agg_df_formatted = agg_df.copy()
numeric_columns = ["avg_time", "min_time", "max_time", "cpu_avg", "cpu_min", "cpu_max", "memory_avg", "memory_min", "memory_max"]

for col in numeric_columns:
    agg_df_formatted[col] = agg_df_formatted[col].round(2)

# === Создаём сетку подграфов: таблица + гистограмма ===
specs = [[{"type": "domain"}, {"type": "xy"}] for _ in subdirs]

fig = make_subplots(
    rows=len(subdirs),
    cols=2,
    specs=specs,
    subplot_titles=[
        f"{cat} — усреднённые метрики" if i % 2 == 0 else f"{cat} — гистограмма времени"
        for cat in subdirs for i in range(2)
    ],
    horizontal_spacing=0.12,
    vertical_spacing=0.15,
)

# === Добавляем данные по каждой категории ===
for i, cat in enumerate(subdirs):
    cat_data = df[df["category"] == cat]
    cat_agg = agg_df[agg_df["category"] == cat].iloc[0]

    # --- Таблица средних значений ---
    table_data = [
        ["Показатель", "Min", "Max", "Avg"],
        ["CPU (%)", f"{cat_agg.cpu_min:.1f}", f"{cat_agg.cpu_max:.1f}", f"{cat_agg.cpu_avg:.1f}"],
        ["Память (МБ)", f"{cat_agg.memory_min:.1f}", f"{cat_agg.memory_max:.1f}", f"{cat_agg.memory_avg:.1f}"],
        ["Время (с)", f"{cat_agg.min_time:.2f}", f"{cat_agg.max_time:.2f}", f"{cat_agg.avg_time:.2f}"]
    ]

    fig.add_trace(
        go.Table(
            header=dict(values=table_data[0], fill_color="lightgrey", align="center"),
            cells=dict(values=list(zip(*table_data[1:])), align="center"),
        ),
        row=i + 1, col=1
    )

    # --- Гистограмма времени выполнения ---
    exec_times = sum(cat_data["execution_times"].tolist(), [])
    fig.add_trace(
        go.Histogram(
            x=exec_times, 
            nbinsx=10, 
            marker_color="cornflowerblue",
            name="Время выполнения",
            hovertemplate="<b>Диапазон:</b> %{x} сек<br><b>Количество:</b> %{y}<extra></extra>"
        ),
        row=i + 1, col=2
    )
    
    # Добавляем подписи осей для гистограмм
    fig.update_xaxes(
        title_text="Время выполнения (секунды)",
        row=i + 1, col=2
    )
    fig.update_yaxes(
        title_text="Количество измерений",
        row=i + 1, col=2
    )

# === Общий макет ===
fig.update_layout(
    height=320 * len(subdirs),
    width=1200,
    title_text="📊 Анализ результатов тестов по категориям",
    showlegend=False
)

# === Сводная таблица с отформатированными данными ===
summary_fig = go.Figure(
    data=[
        go.Table(
            header=dict(values=list(agg_df_formatted.columns), fill_color="lightgrey", align="center"),
            cells=dict(values=[agg_df_formatted[col] for col in agg_df_formatted.columns], align="center")
        )
    ]
)
summary_fig.update_layout(title="Сводная таблица сравнения категорий")

# === Сохраняем всё в один HTML ===
with open(output_html, "w") as f:
    f.write(fig.to_html(full_html=False, include_plotlyjs='cdn'))
    f.write("<hr><h2 style='text-align:center;'>Сводная таблица сравнения категорий</h2>")
    f.write(summary_fig.to_html(full_html=False, include_plotlyjs=False))

print(f"✅ Отчёт сохранён: {output_html}")
import pandas as pd
import numpy as np

pd.set_option("display.max_rows", None)

NON_DATE_VALUES = {
    "deadline not specified",
    "on course",
    "vencido",
    "the deadline for applying has passed",
    "ongoing",
    "deadline needs updating",
    "currently suspended",
}

MONTHS_ES = {
    "enero": "january",
    "febrero": "february",
    "marzo": "march",
    "abril": "april",
    "mayo": "may",
    "junio": "june",
    "julio": "july",
    "agosto": "august",
    "septiembre": "september",
    "setiembre": "september",
    "octubre": "october",
    "noviembre": "november",
    "diciembre": "december",
    "dic": "dec",
}

dataset_raw = pd.read_excel("resources/planillas/raw_funding_spreadsheet.xlsx")
dataset_raw = dataset_raw.drop(dataset_raw.columns[0], axis=1)

relevant_data = dataset_raw.drop(columns=[
    "Small Grant or particularly relevant",
    "Checked",
    "Publicar en boletin",
    "Publicar en DB Publica",
    "Estatus",
    "Days left",
    "area 3",
    "area 4",
    "area 5",
    "area 6",
    "area 7",
    "nota",
    "Restricciones para aplicar / Elegibilidad",
    "Scope geografico",
    "Scope econo-politico",
    "Scope temático",
    "Destinatarios1",
    "Destinatarios2",
    "Destinatarios3",
    "Destinatarios4",
    "Duracion de los proyectos/ Periodo de financiación",
    "periodicidad",
    "dead 1",
    "dead 2",
    "dead 3",
    "dead 4",
    "Fecha Decision/Publicación de resultados",
    "Aplicar mas de una vez",
    "Aim",
    "N",
    "source (scraped from. If empty manually scraped )",
    "keywords queried",
    "ofibusubID",
    "Unnamed: 48",
    "Unnamed: 49",
    "Unnamed: 50",
    "Unnamed: 51",
    "Merged Doc ID - BOLETIN N°16 DICIEMBRE 2023",
    "Merged Doc URL - BOLETIN N°16 DICIEMBRE 2023",
    "Link to merged Doc - BOLETIN N°16 DICIEMBRE 2023",
    "Document Merge Status - BOLETIN N°16 DICIEMBRE 2023",
    "Merged Doc ID - BOLETIN MARZO 2024",
    "Merged Doc URL - BOLETIN MARZO 2024",
    "Link to merged Doc - BOLETIN MARZO 2024",
    "Document Merge Status - BOLETIN MARZO 2024"
], errors="ignore")

dataset = relevant_data.rename(columns={
    "Nombre": "name",
    "Tipo": "type",
    "Gran area 1": "main_area",
    "Gran area 2": "secondary_area",
    "Based on": "based_on",
    "Entidad otorgante": "grantor",
    "Descripcion": "description",
    "Link": "link",
    "Fecha próximo deadline": "deadline",
    "Monto en moneda original": "amount_original",
    "Moneda original": "currency"
})

clean = dataset.copy()

clean = clean[
    clean["name"].notna() &
    clean["name"].astype(str).str.strip().ne("")
]

clean = clean[
    clean["type"].notna() &
    clean["type"].astype(str).str.strip().ne("") &
    clean["type"].astype(str).str.strip().ne("titulo")
]

clean = clean[
    clean["main_area"].notna() &
    clean["main_area"].astype(str).str.strip().ne("")
]

clean["main_area"] = clean["main_area"].replace({
    "all areas of interest": "ALL AREA OF INTEREST",
    "Engineering AND \nARCHITECTURE": "Engineering AND ARCHITECTURE",
    "ALL AREA PF INTEREST": "ALL AREA OF INTEREST",
    "Public health data systems.": "Public health data systems"
})

clean = clean[
    clean["deadline"].notna() &
    clean["deadline"].astype(str).str.strip().ne("") &
    clean["deadline"].astype(str).str.strip().ne("#REF!")
]

clean["deadline"] = clean["deadline"].replace({
    "an 02, 2024": "2024-01-02"
})

clean["deadline"] = clean["deadline"].apply(
    lambda x: x.split(";")[0] if ";" in str(x) else x
)

clean["deadline"] = (
    clean["deadline"]
    .astype(str)
    .str.lower()
    .str.strip()
    .str.replace(r"\.$", "", regex=True)
)

clean["deadline_raw"] = (
    clean["deadline"]
    .str.replace(r"\bde\b", "", regex=True)
    .str.replace(r"\s+", " ", regex=True)
    .str.strip()
)

clean.loc[
    clean["deadline_raw"].isin(NON_DATE_VALUES),
    "deadline_raw"
] = np.nan

for es, en in MONTHS_ES.items():
    clean["deadline_raw"] = clean["deadline_raw"].str.replace(es, en, regex=False)

clean["deadline_clean"] = pd.to_datetime(
    clean["deadline_raw"],
    errors="coerce",
    dayfirst=True
)

clean["deadline"] = np.where(
    clean["deadline_clean"].notna(),
    clean["deadline_clean"].dt.strftime("%Y-%m-%d"),
    clean["deadline"]
)

clean = clean.drop(columns=["deadline_clean", "deadline_raw"])

clean["currency"] = np.where(
    clean["amount_original"].notna() & clean["amount_original"].astype(str).str.strip().ne(""),
    clean["currency"],
    "USD"
)

clean["amount"] = np.where(
    clean["amount_original"].notna() & clean["amount_original"].astype(str).str.strip().ne(""),
    clean["amount_original"].astype(str),
    np.where(
        clean["monto minimo (aprox en U$S)"].notna() & clean["monto maximo (aprox en U$S)"].notna(),
        clean["monto minimo (aprox en U$S)"].astype(str) + " - " +
        clean["monto maximo (aprox en U$S)"].astype(str),
        clean["monto tipico/promedio otorgado (U$S)"].astype(str)
    )
)

# Se intercambian los "nan" por "not-specified", que tomará el traductor para
# traducirlo al idioma adecuado.
clean["amount"] = (
    clean["amount"]
    .replace(
        ["nan", "None", "Not specified", "Not Specified"],
        "not-specified"
    )
    .fillna("not-specified")
)

# Se establece el dólar como la "currency" por defecto.
clean["currency"] = (
    clean["currency"]
    .replace(["nan", "None"], np.nan)
    .fillna("USD")
)

clean = clean.drop(columns=[
    "amount_original",
    "monto minimo (aprox en U$S)",
    "monto maximo (aprox en U$S)",
    "monto tipico/promedio otorgado (U$S)"
], errors="ignore")

# Reordenamiento de columnas.
cols = list(clean.columns)
cols.remove("amount")
cols.insert(cols.index("currency") + 1, "amount")
clean = clean[cols]

clean.to_excel("resources/planillas/clean_funding_spreadsheet.xlsx", index=False)

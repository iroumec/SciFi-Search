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
    "the deadline for applying has passed",
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

dataset_raw = pd.read_excel("resources/planillas/funding.xlsx")

#Se descarta la columna "2022-08-02 09:23:37.436000"
dataset_raw = dataset_raw.drop(dataset_raw.columns[0], axis=1)

dataset_raw["Monto en moneda original"].value_counts(dropna=False)
dataset_raw["monto minimo (aprox en U$S)"].value_counts(dropna=False)
dataset_raw["monto maximo (aprox en U$S)"].value_counts(dropna=False)

#Se descartan las columnas que no interesan.
relevant_data = dataset_raw.drop(columns=[
                                    #"2022-08-02 09:23:37.436000", 
                                    "Small Grant or particularly relevant", 
                                    "Checked",
                                    "Publicar en boletin",
                                    "Publicar en DB Publica",
                                    "Estatus",
                                    "Days left",
                                    #"Nombre",
                                    #"Tipo",
                                    #"Gran area 1",
                                    #"Gran area 2",
                                    "area 3",
                                    "area 4",
                                    "area 5",
                                    "area 6",
                                    "area 7",
                                    #"Link",
                                    #"Descripcion",
                                    "nota",
                                    "Restricciones para aplicar / Elegibilidad",
                                    #"Based on",
                                    #"Entidad otorgante",
                                    "Scope geografico",
                                    "Scope econo-politico",
                                    "Scope temático",
                                    "Destinatarios1",
                                    "Destinatarios2",
                                    "Destinatarios3",
                                    "Destinatarios4",
                                    "Duracion de los proyectos/ Periodo de financiación",
                                    #"Moneda original",
                                    #"Monto en moneda original",
                                    "monto minimo (aprox en U$S)",
                                    "monto maximo (aprox en U$S)",
                                    "monto tipico/promedio otorgado (U$S)",
                                    "periodicidad",
                                    "dead 1",
                                    "dead 2",
                                    "dead 3",
                                    "dead 4",
                                    #"Fecha próximo deadline",
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
    "Entidad otorgante": "grantor",
    "Descripcion": "description",
    "Link": "link",
    "Fecha próximo deadline": "deadline",
    "Monto en moneda original": "amount_original",
    "Moneda original": "currency"
})


clean = dataset.copy()

#Se descartan las filas con nombre nulo o vacíos
clean = clean[
    clean["name"].notna() &                         
    clean["name"].astype(str).str.strip().ne("")    
]

#Se descartan las filas con tipo nulo, tipo inválido o tipo vacío
clean = clean[
    clean["type"].notna() &                             
    clean["type"].astype(str).str.strip().ne("") &      
    clean["type"].astype(str).str.strip().ne("titulo")  
]

#Se descartan las filas con valor de gran aréa 1 nulo o vacío
clean = clean[
    clean["main_area"].notna() &                         
    clean["main_area"].astype(str).str.strip().ne("")    
]

#Se corrigen valores específicos de la columna
clean["main_area"] = clean["main_area"].replace({
    "all areas of interest": "ALL AREA OF INTEREST",
    "Engineering AND \nARCHITECTURE": "Engineering AND ARCHITECTURE",
    "ALL AREA PF INTEREST": "ALL AREA OF INTEREST",
    "Public health data systems.": "Public health data systems"
})

#Se descartan las filas con valor de deadline nulo o vacío
clean = clean[
    clean["deadline"].notna() &                         
    clean["deadline"].astype(str).str.strip().ne("") &
    clean["deadline"].astype(str).str.strip().ne("#REF!")   
]

#Se corrigen valores específicos
clean["deadline"] = clean["deadline"].replace({
    "an 02, 2024": "2024-01-02"
})

#Se corrigen duplicados separados por ;
clean["deadline"] = clean["deadline"].apply(
    lambda x: x.split(";")[0] if ";" in x else x
)

#Se normalizan valores
clean["deadline"] = (
    clean["deadline"]
    .astype(str)
    .str.strip()
    .str.lower()
    .str.replace(r"\.$", "", regex=True)
    .str.strip()
)

#Se crea una nueva columna con valores en proceso de ser normalizados
clean["deadline_raw"] = (
    clean["deadline"]
    .str.replace(r"\bde\b", "", regex=True)
    .str.replace(r"\s+", " ", regex=True)
    .str.strip()
)

#Reemplaza valores que no son fechas con NaN
clean.loc[
    clean["deadline_raw"].isin(NON_DATE_VALUES),
    "deadline_raw"
] = np.nan

#Traduccion de español a inglés
for es, en in MONTHS_ES.items():
    clean["deadline_raw"] = clean["deadline_raw"].str.replace(
        es, en, regex=False
    )

#Normalización de valores a fechas (los NaN se vuelven NaT)
clean["deadline_clean"] = pd.to_datetime(
    clean["deadline_raw"],
    errors="coerce",
    dayfirst=True
)

#Unificación
clean["deadline"] = np.where(
    clean["deadline_clean"].notna(),
    clean["deadline_clean"].dt.strftime("%Y-%m-%d"),
    clean["deadline"].astype(str)
)

#Se descartan las columnas auxiliares
clean = clean.drop(columns=["deadline_clean","deadline_raw"])

## MONTO ##
##

clean.to_excel("resources/planillas/clean_funding.xlsx", index=False)
import pandas as pd

# Lectura del Excel
df = pd.read_excel("resources/planillas/funding.xlsx")

# Conversión a JSON sin orientación de registros
df.to_json("resources/planillas/funding.json")

# Conversión a JSON con orientación de registros
df.to_json("resources/planillas/fundingRecords.json", orient="records", lines=False, date_format="iso")

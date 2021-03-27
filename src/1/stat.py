import pandas as pd

# https://realpython.com/pandas-python-explore-dataset/

# Used code from: 
# https://stackoverflow.com/questions/21269399/datetime-dtypes-in-pandas-read-csv

df = pd.read_csv(r'./data/Sociologove.csv', parse_dates = ['born', 'died'])

print(df.to_string(), df.dtypes)
s = df['born']
print(s.where(s < 1960))

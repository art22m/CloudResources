import requests

url = "https://mts-olimp-cloud.codenrock.com/api/resource?token=some-token"
resp = requests.get(url).json()

for entity in resp:
    print(entity)




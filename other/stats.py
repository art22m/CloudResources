import requests

url_stats = "https://mts-olimp-cloud.codenrock.com/api/statistic?token=some-token"
resp_stats = requests.get(url_stats).json()

url_resources = "https://mts-olimp-cloud.codenrock.com/api/resource?token=some-token"
resp_resources = requests.get(url_resources).json()

# url_prices = "https://mts-olimp-cloud.codenrock.com/api/price"
# resp_prices = requests.get(url_prices).json()

db_count_work = 0
db_count_failed = 0
vm_count_work = 0
vm_count_failed = 0
for entity in resp_resources:
    if entity["type"] == "db":
        if entity["failed"]:
            db_count_failed += 1
        else:
            db_count_work += 1
    elif entity["type"] == "vm":
        if entity["failed"]:
            vm_count_failed += 1
        else:
            vm_count_work += 1
    else:
        print("unknown type")

print("timestamp", resp_stats['timestamp'])
print("cost_total", resp_stats['cost_total'])
print("response_time", resp_stats['response_time'])
print("last_hour", resp_stats['lastHour'])
print("availability", resp_stats['availability'])

print("online", resp_stats['online'])
print("requests", resp_stats['requests'])
print("--------------")
print("db_count_work", db_count_work)
print("db_count_failed", db_count_failed)
print("db_cpu", resp_stats['db_cpu'])
print("db_cpu_load", resp_stats['db_cpu_load'])
print("db_ram", resp_stats['db_ram'])
print("db_ram_load", resp_stats['db_ram_load'])
print("--------------")
print("vm_count_work", vm_count_work)
print("vm_count_failed", vm_count_failed)
print("vm_cpu", resp_stats['vm_cpu'])
print("vm_cpu_load", resp_stats['vm_cpu_load'])
print("vm_ram", resp_stats['vm_ram'])
print("vm_ram_load", resp_stats['vm_ram_load'])

print("-----PRICES-----")
# print(resp_prices)

import os
import json

def process_results_line(l):
    l = json.loads(l)
    name = l["name"]
    value = (l["measures"])["value"]
    item = {}
    item["name"] = name
    item["value"] = value
    return item

#####################################################
## Merges all results into one big list of metrics ##
#####################################################
def aggregate_results(results_dir):
    res = []
    for subdir, _, files in os.walk(results_dir):
        for filename in files:
            filepath = subdir + os.sep + filename
            if filepath.split("/")[-1] == "results.out":
                resultFile = open(filepath, 'r')
                for l in resultFile.readlines():
                    res.append(process_results_line(l))
    return res, len(os.listdir(results_dir))   

###############################################################################
## Processes all metrics produced in FastDelivery and compiles the means and ##
## max for each run and averages it among the different executed runs        ##
###############################################################################
def digested_results(res, test):
    
    summary = {}
    final_data = {}
    runs = len(res)/len(test)
    
    for measure in test:
        summary[measure] = 0
        final_data[measure + "/max"] = 0
    for item in res:
        if item["name"] in test:
            summary[item["name"]] += item["value"]
            if item["value"] > final_data[item["name"] + "/max"]:
                final_data[item["name"] + "/max"] = item["value"]

    for sums in summary:
        final_data[sums + "/mean"] = summary[sums]/runs

    return final_data

##################
## Main Program ##
##################
metrics = {"fast":['Avg time to sub - FastDelivery', 'Events received - FastDelivery', 'Avg latency of events - FastDelivery',
     'Avg time to sub - FastDelivery', 'CPU used - FastDelivery', 'Memory used - FastDelivery']}
dir_path = os.path.dirname(os.path.realpath(__file__))
agg, testcases = aggregate_results(dir_path + "/../../../data/outputs/local_docker/contentps-test/")
final_res = digested_results(agg, metrics["fast"])
print(final_res)
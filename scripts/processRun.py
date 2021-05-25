import os
import json
import sys

########################################################
## Interprets a line (metric) from a results.out file ##
########################################################
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
    num_metrics = 0
    
    for measure in test:
        summary[measure] = 0
        final_data[measure + "/max"] = 0
    for item in res:
        if item["name"] in test:
            num_metrics += 1
            summary[item["name"]] += item["value"]
            if item["value"] > final_data[item["name"] + "/max"]:
                final_data[item["name"] + "/max"] = item["value"]

    peer_runs = num_metrics/len(test)
    for sums in summary:
        final_data[sums + "/mean"] = summary[sums]/peer_runs

    final_data["runs"] = peer_runs/18
    return final_data

##################
## Main Program ##
##################
metrics = {"fast": ['Avg time to sub - FastDelivery', 'Avg latency of events - FastDelivery','Avg time to sub - FastDelivery', 
'CPU used - FastDelivery', 'Memory used - FastDelivery'],
 "normal_scout": ['Avg event latency - ScoutSubs normal', 'Avg time to sub - ScoutSubs normal',
'CPU used - ScoutSubs normal', 'Memory used - ScoutSubs normal'],
 "normal_subBurst": ['Avg event latency - ScoutSubs subBurst', 'Avg time to sub - ScoutSubs subBurst', 
'CPU used - ScoutSubs subBurst', 'Memory used - ScoutSubs subBurst'],
 "normal_eventBurst": ['Avg event latency - ScoutSubs eventBurst', 'Avg time to sub - ScoutSubs eventBurst',
'CPU used - ScoutSubs eventBurst', 'Memory used - ScoutSubs eventBurst']}

dir_path = os.path.dirname(os.path.realpath(__file__))
agg, testcases = aggregate_results(dir_path + "/../../../data/outputs/local_docker/contentps-test/")
final_res = digested_results(agg, metrics["fast"])
for i in final_res:
    print(i + " >> {0}".format(final_res[i]))
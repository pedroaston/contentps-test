import os
import json
import matplotlib.pyplot as plt
import numpy as np

##############################
## Auxiliar data structures ##
##############################
metric_types = {"fast": ['Avg time to sub - FastDelivery', 'Avg event latency - FastDelivery', 
'CPU used - FastDelivery', 'Memory used - FastDelivery'],
 "normal_scout_BU": ['Avg event latency - ScoutSubs normalBU',
'CPU used - ScoutSubs normalBU', 'Memory used - ScoutSubs normalBU'],
 "normal_scout_RU": ['Avg event latency - ScoutSubs normalRU',
'CPU used - ScoutSubs normalRU', 'Memory used - ScoutSubs normalRU'],
 "normal_scout_BR": ['Avg event latency - ScoutSubs normalBR', 'Avg time to sub - ScoutSubs normalBR',
'CPU used - ScoutSubs normalBR', 'Memory used - ScoutSubs normalBR'],
 "normal_scout_RR": ['Avg event latency - ScoutSubs normalRR', 'Avg time to sub - ScoutSubs normalRR',
'CPU used - ScoutSubs normalRR', 'Memory used - ScoutSubs normalRR'],
 "normal_subBurst": ['Avg event latency - ScoutSubs subBurst', 'Avg time to sub - ScoutSubs subBurst', 
'CPU used - ScoutSubs subBurst', 'Memory used - ScoutSubs subBurst'],
 "normal_eventBurst": ['Avg event latency - ScoutSubs eventBurst', 'Avg time to sub - ScoutSubs eventBurst',
'CPU used - ScoutSubs eventBurst', 'Memory used - ScoutSubs eventBurst']}

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
    return res   

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
        final_data[sums + "/mean"] = round(summary[sums]/peer_runs, 2)

    final_data["runs"] = peer_runs//18
    return final_data

###############################################
## Returns summary of the interested metrics ##
###############################################
def metric_summary(type):
    dir_path = os.path.dirname(os.path.realpath(__file__))
    agg = aggregate_results(dir_path + "/../../../data/outputs/local_docker/contentps-test/")
    final_res = digested_results(agg, metric_types[type])
    return final_res

#########################
## Memory metrics plot ##
#########################
def plot_memory_metric(scenario):
    fast_res = metric_summary("fast")
    scout_res_BU = metric_summary(scenario + "_scout_BU")
    scout_res_BR = metric_summary(scenario + "_scout_BR")
    scout_res_RU = metric_summary(scenario + "_scout_RU")
    scout_res_RR = metric_summary(scenario + "_scout_RR")

    labels = ['FastDelivery', 'Base-Unreliable', 'Base-Reliable', 'Redirect-Unreliable', 'Redirect-Reliable']
    mean_values = [fast_res['Memory used - FastDelivery/mean'], scout_res_BU['Memory used - ScoutSubs '+scenario+'BU/mean'],
     scout_res_BR['Memory used - ScoutSubs '+scenario+'BR/mean'], scout_res_RU['Memory used - ScoutSubs '+scenario+'RU/mean'],
     scout_res_RR['Memory used - ScoutSubs '+scenario+'RR/mean']]
    max_values = [fast_res['Memory used - FastDelivery/max'],scout_res_BU['Memory used - ScoutSubs '+scenario+'BU/max'],
     scout_res_BR['Memory used - ScoutSubs '+scenario+'BR/max'], scout_res_RU['Memory used - ScoutSubs '+scenario+'RU/max'],
     scout_res_RR['Memory used - ScoutSubs '+scenario+'RR/max']]

    x = np.arange(len(labels))  # the label locations
    width = 0.4  # the width of the bars

    fig, ax = plt.subplots()
    rects1 = ax.bar(x - width/2, mean_values, width, label='mean')
    rects2 = ax.bar(x + width/2, max_values, width, label='max')

    # Add some text for labels, title and custom x-axis tick labels, etc.
    ax.set_ylabel('bytes used')
    ax.set_xlabel('variants')
    ax.set_title('Memory used by pubsub')
    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend()

    ax.bar_label(rects1, padding=3)
    ax.bar_label(rects2, padding=3)

    fig.tight_layout()

    plt.show()

######################
## Cpu metrics plot ##
######################
def plot_cpu_metric(scenario):
    fast_res = metric_summary("fast")
    scout_res_BU = metric_summary(scenario + "_scout_BU")
    scout_res_BR = metric_summary(scenario + "_scout_BR")
    scout_res_RU = metric_summary(scenario + "_scout_RU")
    scout_res_RR = metric_summary(scenario + "_scout_RR")

    labels = ['FastDelivery', 'Base-Unreliable', 'Base-Reliable', 'Redirect-Unreliable', 'Redirect-Reliable']
    mean_values = [fast_res['CPU used - FastDelivery/mean'], scout_res_BU['CPU used - ScoutSubs '+scenario+'BU/mean'],
     scout_res_BR['CPU used - ScoutSubs '+scenario+'BR/mean'], scout_res_RU['CPU used - ScoutSubs '+scenario+'RU/mean'],
     scout_res_RR['CPU used - ScoutSubs '+scenario+'RR/mean']]
    max_values = [fast_res['CPU used - FastDelivery/max'],scout_res_BU['CPU used - ScoutSubs '+scenario+'BU/max'],
     scout_res_BR['CPU used - ScoutSubs '+scenario+'BR/max'], scout_res_RU['CPU used - ScoutSubs '+scenario+'RU/max'],
     scout_res_RR['CPU used - ScoutSubs '+scenario+'RR/max']]

    x = np.arange(len(labels))  # the label locations
    width = 0.4  # the width of the bars

    fig, ax = plt.subplots()
    rects1 = ax.bar(x - width/2, mean_values, width, label='mean')
    rects2 = ax.bar(x + width/2, max_values, width, label='max')

    # Add some text for labels, title and custom x-axis tick labels, etc.
    ax.set_ylabel('cpu user-time (s)')
    ax.set_xlabel('variants')
    ax.set_title('CPU time used by pubsub')
    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend()

    ax.bar_label(rects1, padding=3)
    ax.bar_label(rects2, padding=3)

    fig.tight_layout()

    plt.show()

####################################
## Avg event latency metrics plot ##
####################################
def plot_latency_metric(scenario):
    fast_res = metric_summary("fast")
    scout_res_BU = metric_summary(scenario + "_scout_BU")
    scout_res_BR = metric_summary(scenario + "_scout_BR")
    scout_res_RU = metric_summary(scenario + "_scout_RU")
    scout_res_RR = metric_summary(scenario + "_scout_RR")

    labels = ['FastDelivery', 'Base-Unreliable', 'Base-Reliable', 'Redirect-Unreliable', 'Redirect-Reliable']
    mean_values = [fast_res['Avg event latency - FastDelivery/mean'], scout_res_BU['Avg event latency - ScoutSubs '+scenario+'BU/mean'],
     scout_res_BR['Avg event latency - ScoutSubs '+scenario+'BR/mean'], scout_res_RU['Avg event latency - ScoutSubs '+scenario+'RU/mean'],
     scout_res_RR['Avg event latency - ScoutSubs '+scenario+'RR/mean']]
    max_values = [fast_res['Avg event latency - FastDelivery/max'],scout_res_BU['Avg event latency - ScoutSubs '+scenario+'BU/max'],
     scout_res_BR['Avg event latency - ScoutSubs '+scenario+'BR/max'], scout_res_RU['Avg event latency - ScoutSubs '+scenario+'RU/max'],
     scout_res_RR['Avg event latency - ScoutSubs '+scenario+'RR/max']]

    x = np.arange(len(labels))  # the label locations
    width = 0.4  # the width of the bars

    fig, ax = plt.subplots()
    rects1 = ax.bar(x - width/2, mean_values, width, label='mean')
    rects2 = ax.bar(x + width/2, max_values, width, label='max')

    # Add some text for labels, title and custom x-axis tick labels, etc.
    ax.set_ylabel('event latency (ms)')
    ax.set_xlabel('variants')
    ax.set_title('Event latency with pubsub')
    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend()

    ax.bar_label(rects1, padding=3)
    ax.bar_label(rects2, padding=3)

    fig.tight_layout()

    plt.show()
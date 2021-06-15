import os
import json
import matplotlib.pyplot as plt
import numpy as np
import seaborn as sns

##############################
## Auxiliar data structures ##
##############################
avgMetrics = ['CPU used', 'Memory used', 'Event Latency', 'Sub Latency']
cumMetrics = ['# Events Missing', '# Events Duplicated']

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
                    item = process_results_line(l)
                    if item["value"] >= 0:
                        res.append(item)
    return res

###############################################################################
## Processes all metrics produced in FastDelivery and compiles the means and ##
## max for each run and averages it among the different executed runs        ##
###############################################################################
def digested_results(res, test):

    testAvgMetrics = []
    for m in avgMetrics:
        testAvgMetrics.append("{0} - {1}".format(m, test))
    testCumMetrics = []
    for m in cumMetrics:
        testCumMetrics.append("{0} - {1}".format(m, test))

    summary = {}
    final_data = {}
    num_metrics = {}
    
    for measure in testAvgMetrics:
        summary[measure] = 0
        num_metrics[measure] = 0
        final_data[measure+"/max"] = 0
        final_data[measure+"/all"] = []
    for measure in testCumMetrics:
        final_data[measure] = 0

    for item in res:
        if item["name"] in testAvgMetrics:
            num_metrics[item["name"]] += 1
            summary[item["name"]] += item["value"]
            final_data[item["name"]+"/all"].append(item["value"])
            if final_data[item["name"]+"/max"] < item["value"]:
                final_data[item["name"]+"/max"] = item["value"]
        elif item["name"] in testCumMetrics:
            final_data[item["name"]] += item["value"]

    for sums in summary:
        if num_metrics[sums] != 0:
            final_data[sums + "/mean"] = summary[sums]/num_metrics[sums]
        else :
            final_data[sums + "/mean"] = 0

    return final_data

###############################################
## Returns summary of the interested metrics ##
###############################################
def metric_summary(type):
    dir_path = os.path.dirname(os.path.realpath(__file__))
    agg = aggregate_results(dir_path + "/../../../data/outputs/local_docker/contentps-test/")
    final_res = digested_results(agg, type)
    return final_res

##############################################
## Shows all plots of a particular scenario ##
##############################################
def plot_everything(scenario):
    fast = metric_summary("FastDelivery")
    scout_BU = metric_summary("{0} {1}{2}".format("ScoutSubs", scenario,"BU"))
    scout_BR = metric_summary("{0} {1}{2}".format("ScoutSubs", scenario,"BR"))
    scout_RU = metric_summary("{0} {1}{2}".format("ScoutSubs", scenario,"RU"))
    scout_RR = metric_summary("{0} {1}{2}".format("ScoutSubs", scenario,"RR"))

    plot_correctness_metrics(fast, scout_BU, scout_BR, scout_RU, scout_RR, scenario)
    plot_latency_metric(fast, scout_BU, scout_BR, scout_RU, scout_RR, scenario)
    boxplot_event_latency(fast, scout_BU, scout_BR, scout_RU, scout_RR, scenario)
    boxplot_sub_latency(fast, scout_BR, scout_RR, scenario)
    plot_memory_metric(fast, scout_BU, scout_BR, scout_RU, scout_RR, scenario)
    plot_cpu_metric(fast, scout_BU, scout_BR, scout_RU, scout_RR, scenario)

###########################
## Event Latency boxplot ##
###########################
def boxplot_event_latency(fast, scout_BU, scout_BR, scout_RU, scout_RR, scenario):

    data = [fast['Event Latency - FastDelivery/all'], scout_BU['Event Latency - ScoutSubs '+scenario+'BU/all'], 
     scout_BR['Event Latency - ScoutSubs '+scenario+'BR/all'], scout_RU['Event Latency - ScoutSubs '+scenario+'RU/all'],
     scout_RR['Event Latency - ScoutSubs '+scenario+'RR/all']]

    labels = ['FastDelivery', 'Basic-Unreliable', 'Basic-Reliable', 'Redirect-Unreliable', 'Redirect-Reliable']

    sns.set_context('talk', font_scale = 0.75)
    fig7, ax7 = plt.subplots(figsize=(10, 8))
    ax7.set_title('Event Latency Distribution', pad=30, fontsize=20)
    ax7.boxplot(data, labels=labels, patch_artist=True)
    ax7.set_xlabel('Variants', labelpad=20)
    ax7.set_ylabel('Time (ms)', labelpad=20)

    plt.show()

#########################
## Sub Latency boxplot ##
#########################
def boxplot_sub_latency(fast, scout_BR, scout_RR, scenario):

    data = [fast['Sub Latency - FastDelivery/all'], scout_BR['Sub Latency - ScoutSubs '+scenario+'BR/all'],
     scout_RR['Sub Latency - ScoutSubs '+scenario+'RR/all']]

    labels = ['FastDelivery', 'Basic-Reliable', 'Redirect-Reliable']

    sns.set_context('talk', font_scale = 0.75)
    fig7, ax7 = plt.subplots(figsize=(10, 8))
    ax7.set_title('Sub Latency Distribution', pad=30, fontsize=20)
    ax7.boxplot(data, labels=labels, patch_artist=True)
    ax7.set_xlabel('Variants', labelpad=20)
    ax7.set_ylabel('Time (ms)', labelpad=20)

    plt.show()

#########################
## Memory metrics plot ##
#########################
def plot_memory_metric(fast, scout_BU, scout_BR, scout_RU, scout_RR, scenario):

    labels = ['FastDelivery', 'Base-Unreliable', 'Base-Reliable', 'Redirect-Unreliable', 'Redirect-Reliable']
    mean_values = [fast['Memory used - FastDelivery/mean'], scout_BU['Memory used - ScoutSubs '+scenario+'BU/mean'],
     scout_BR['Memory used - ScoutSubs '+scenario+'BR/mean'], scout_RU['Memory used - ScoutSubs '+scenario+'RU/mean'],
     scout_RR['Memory used - ScoutSubs '+scenario+'RR/mean']]
    max_values = [fast['Memory used - FastDelivery/max'],scout_BU['Memory used - ScoutSubs '+scenario+'BU/max'],
     scout_BR['Memory used - ScoutSubs '+scenario+'BR/max'], scout_RU['Memory used - ScoutSubs '+scenario+'RU/max'],
     scout_RR['Memory used - ScoutSubs '+scenario+'RR/max']]


    sns.set_context('talk', font_scale = 0.75)
    fig, ax = plt.subplots(figsize=(12, 8))

    x = np.arange(len(labels))  
    width = 0.4  

    
    ax.bar(x - width/2, mean_values, width, label='mean')
    ax.bar(x + width/2, max_values, width, label='max')

    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend()

    ax.spines['top'].set_visible(False)
    ax.spines['right'].set_visible(False)
    ax.spines['left'].set_visible(False)
    ax.spines['bottom'].set_color('#DDDDDD')
    ax.tick_params(bottom=False, left=False)
    ax.set_axisbelow(True)
    ax.yaxis.grid(True, color='#EEEEEE')
    ax.xaxis.grid(False)

    ax.set_ylabel('# MB used', labelpad=20)
    ax.set_xlabel('Variants', labelpad=20)
    ax.set_title('Memory used by pubsub', pad=30, fontsize=20)

    for bar in ax.patches:
        bar_value = bar.get_height()
        text = f'{bar_value:.1f}'
        text_x = bar.get_x() + bar.get_width() / 2
        text_y = bar.get_y() + bar_value
        bar_color = bar.get_facecolor()
        ax.text(text_x, text_y, text, ha='center', va='bottom', color=bar_color, size=12)

    fig.tight_layout()
    plt.show()
    print()

######################
## Cpu metrics plot ##
######################
def plot_cpu_metric(fast, scout_BU, scout_BR, scout_RU, scout_RR, scenario):

    labels = ['FastDelivery', 'Base-Unreliable', 'Base-Reliable', 'Redirect-Unreliable', 'Redirect-Reliable']
    mean_values = [fast['CPU used - FastDelivery/mean'], scout_BU['CPU used - ScoutSubs '+scenario+'BU/mean'],
     scout_BR['CPU used - ScoutSubs '+scenario+'BR/mean'], scout_RU['CPU used - ScoutSubs '+scenario+'RU/mean'],
     scout_RR['CPU used - ScoutSubs '+scenario+'RR/mean']]
    max_values = [fast['CPU used - FastDelivery/max'],scout_BU['CPU used - ScoutSubs '+scenario+'BU/max'],
     scout_BR['CPU used - ScoutSubs '+scenario+'BR/max'], scout_RU['CPU used - ScoutSubs '+scenario+'RU/max'],
     scout_RR['CPU used - ScoutSubs '+scenario+'RR/max']]

    sns.set_context('talk', font_scale = 0.75)
    fig, ax = plt.subplots(figsize=(12, 8))
    x = np.arange(len(labels))  
    width = 0.4  

    ax.bar(x - width/2, mean_values, width, label='mean')
    ax.bar(x + width/2, max_values, width, label='max')

    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend()

    ax.spines['top'].set_visible(False)
    ax.spines['right'].set_visible(False)
    ax.spines['left'].set_visible(False)
    ax.spines['bottom'].set_color('#DDDDDD')
    ax.tick_params(bottom=False, left=False)
    ax.set_axisbelow(True)
    ax.yaxis.grid(True, color='#EEEEEE')
    ax.xaxis.grid(False)

    ax.set_ylabel('cpu user-time (s)', labelpad=20)
    ax.set_xlabel('Variants', labelpad=20)
    ax.set_title('CPU time used by pubsub', pad=30, fontsize=20)

    for bar in ax.patches:
        bar_value = bar.get_height()
        text = f'{bar_value:.2f}'
        text_x = bar.get_x() + bar.get_width() / 2
        text_y = bar.get_y() + bar_value
        bar_color = bar.get_facecolor()
        ax.text(text_x, text_y, text, ha='center', va='bottom', color=bar_color, size=12)

    fig.tight_layout()
    plt.show()
    print()

####################################
## Avg event latency metrics plot ##
####################################
def plot_latency_metric(fast, scout_BU, scout_BR, scout_RU, scout_RR, scenario):

    labels = ['FastDelivery', 'Base-Unreliable', 'Base-Reliable', 'Redirect-Unreliable', 'Redirect-Reliable']
    mean_values = [fast['Event Latency - FastDelivery/mean'], scout_BU['Event Latency - ScoutSubs '+scenario+'BU/mean'],
     scout_BR['Event Latency - ScoutSubs '+scenario+'BR/mean'], scout_RU['Event Latency - ScoutSubs '+scenario+'RU/mean'],
     scout_RR['Event Latency - ScoutSubs '+scenario+'RR/mean']]
    max_values = [fast['Event Latency - FastDelivery/max'],scout_BU['Event Latency - ScoutSubs '+scenario+'BU/max'],
     scout_BR['Event Latency - ScoutSubs '+scenario+'BR/max'], scout_RU['Event Latency - ScoutSubs '+scenario+'RU/max'],
     scout_RR['Event Latency - ScoutSubs '+scenario+'RR/max']]

    sns.set_context('talk', font_scale = 0.75)
    fig, ax = plt.subplots(figsize=(12, 8))
    x = np.arange(len(labels))
    width = 0.4

    ax.bar(x - width/2, mean_values, width, label='mean')
    ax.bar(x + width/2, max_values, width, label='max')

    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend()

    ax.spines['top'].set_visible(False)
    ax.spines['right'].set_visible(False)
    ax.spines['left'].set_visible(False)
    ax.spines['bottom'].set_color('#DDDDDD')
    ax.tick_params(bottom=False, left=False)
    ax.set_axisbelow(True)
    ax.yaxis.grid(True, color='#EEEEEE')
    ax.xaxis.grid(False)

    ax.set_ylabel('event latency (ms)', labelpad=20)
    ax.set_xlabel('Variants', labelpad=20)
    ax.set_title('Event latency with pubsub', pad=30, fontsize=20)

    for bar in ax.patches:
        bar_value = bar.get_height()
        text = f'{bar_value:.0f}'
        text_x = bar.get_x() + bar.get_width() / 2
        text_y = bar.get_y() + bar_value
        bar_color = bar.get_facecolor()
        ax.text(text_x, text_y, text, ha='center', va='bottom', color=bar_color, size=12)

    fig.tight_layout()
    plt.show()

##############################
## Correctness metrics plot ##
##############################
def plot_correctness_metrics(fast, scout_BU, scout_BR, scout_RU, scout_RR, scenario):

    labels = ['FastDelivery', 'Base-Unreliable', 'Base-Reliable', 'Redirect-Unreliable', 'Redirect-Reliable']
    mean_values = [fast['# Events Missing - FastDelivery'], scout_BU['# Events Missing - ScoutSubs '+scenario+'BU'],
     scout_BR['# Events Missing - ScoutSubs '+scenario+'BR'], scout_RU['# Events Missing - ScoutSubs '+scenario+'RU'],
     scout_RR['# Events Missing - ScoutSubs '+scenario+'RR']]
    max_values = [fast['# Events Duplicated - FastDelivery'],scout_BU['# Events Duplicated - ScoutSubs '+scenario+'BU'],
     scout_BR['# Events Duplicated - ScoutSubs '+scenario+'BR'], scout_RU['# Events Duplicated - ScoutSubs '+scenario+'RU'],
     scout_RR['# Events Duplicated - ScoutSubs '+scenario+'RR']]

    sns.set_context('talk', font_scale = 0.75)
    fig, ax = plt.subplots(figsize=(12, 8))
    x = np.arange(len(labels))
    width = 0.4

    ax.bar(x - width/2, mean_values, width, label='missing')
    ax.bar(x + width/2, max_values, width, label='duplicated')

    ax.set_xticks(x)
    ax.set_xticklabels(labels)
    ax.legend()

    ax.spines['top'].set_visible(False)
    ax.spines['right'].set_visible(False)
    ax.spines['left'].set_visible(False)
    ax.spines['bottom'].set_color('#DDDDDD')
    ax.tick_params(bottom=False, left=False)
    ax.set_axisbelow(True)
    ax.yaxis.grid(True, color='#EEEEEE')
    ax.xaxis.grid(False)

    ax.set_ylabel('# Events', labelpad=20)
    ax.set_xlabel('Variants', labelpad=20)
    ax.set_title('Pubsub Correctness', pad=30, fontsize=20)

    for bar in ax.patches:
        bar_value = bar.get_height()
        text = f'{bar_value:.0f}'
        text_x = bar.get_x() + bar.get_width() / 2
        text_y = bar.get_y() + bar_value
        bar_color = bar.get_facecolor()
        ax.text(text_x, text_y, text, ha='center', va='bottom', color=bar_color, size=12)

    fig.tight_layout()
    plt.show()
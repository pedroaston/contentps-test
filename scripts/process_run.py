import os
import json
import matplotlib.pyplot as plt
import numpy as np
import seaborn as sns

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
 "subBurst_scout_BU": ['Avg event latency - ScoutSubs subBurstBU',  
'CPU used - ScoutSubs subBurstBU', 'Memory used - ScoutSubs subBurstBU'],
 "subBurst_scout_RU": ['Avg event latency - ScoutSubs subBurstRU',  
'CPU used - ScoutSubs subBurstRU', 'Memory used - ScoutSubs subBurstRU'],
 "subBurst_scout_BR": ['Avg event latency - ScoutSubs subBurstBR', 'Avg time to sub - ScoutSubs subBurstBR', 
'CPU used - ScoutSubs subBurstBR', 'Memory used - ScoutSubs subBurstBR'],
 "subBurst_scout_RR": ['Avg event latency - ScoutSubs subBurstRR', 'Avg time to sub - ScoutSubs subBurstRR', 
'CPU used - ScoutSubs subBurstRR', 'Memory used - ScoutSubs subBurstRR'],
 "eventBurst_scout_BU": ['Avg event latency - ScoutSubs eventBurstBU', 
'CPU used - ScoutSubs eventBurstBU', 'Memory used - ScoutSubs eventBurstBU'],
 "eventBurst_scout_RU": ['Avg event latency - ScoutSubs eventBurstRU', 
'CPU used - ScoutSubs eventBurstRU', 'Memory used - ScoutSubs eventBurstRU'],
 "eventBurst_scout_BR": ['Avg event latency - ScoutSubs eventBurstBR', 'Avg time to sub - ScoutSubs eventBurstBR',
'CPU used - ScoutSubs eventBurstBR', 'Memory used - ScoutSubs eventBurstBR'],
 "eventBurst_scout_RR": ['Avg event latency - ScoutSubs eventBurstRR', 'Avg time to sub - ScoutSubs eventBurstRR',
'CPU used - ScoutSubs eventBurstRR', 'Memory used - ScoutSubs eventBurstRR'],
 "fault_scout_BU": ['Avg event latency - ScoutSubs faultBU', 
'CPU used - ScoutSubs faultBU', 'Memory used - ScoutSubs faultBU'],
 "fault_scout_RU": ['Avg event latency - ScoutSubs faultRU', 
'CPU used - ScoutSubs faultRU', 'Memory used - ScoutSubs faultRU'],
 "fault_scout_BR": ['Avg event latency - ScoutSubs faultBR', 'Avg time to sub - ScoutSubs faultBR',
'CPU used - ScoutSubs faultBR', 'Memory used - ScoutSubs faultBR'],
 "fault_scout_RR": ['Avg event latency - ScoutSubs faultRR', 'Avg time to sub - ScoutSubs faultRR',
'CPU used - ScoutSubs faultRR', 'Memory used - ScoutSubs faultRR']}

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
                    if item["value"] == 0:
                        res.append(item)
    return res

###############################################################################
## Processes all metrics produced in FastDelivery and compiles the means and ##
## max for each run and averages it among the different executed runs        ##
###############################################################################
def digested_results(res, test):
    
    summary = {}
    final_data = {}
    num_metrics = {}
    
    for measure in test:
        summary[measure] = 0
        num_metrics[measure] = 0
        final_data[measure + "/max"] = 0
    for item in res:
        if item["name"] in test:
            num_metrics[item["name"]] += 1
            summary[item["name"]] += item["value"]
            if item["value"] > final_data[item["name"] + "/max"]:
                final_data[item["name"] + "/max"] = item["value"]

    for sums in summary:
        if num_metrics[sums] != 0:
            final_data[sums + "/mean"] = summary[sums]/num_metrics[sums]

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


    sns.set_context('talk', font_scale = 0.75)
    fig, ax = plt.subplots(figsize=(12, 8))

    x = np.arange(len(labels))  # the label locations
    width = 0.4  # the width of the bars

    
    rects1 = ax.bar(x - width/2, mean_values, width, label='mean')
    rects2 = ax.bar(x + width/2, max_values, width, label='max')

    # Add some text for labels, title and custom x-axis tick labels, etc.
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

    ax.set_ylabel('# bytes used', labelpad=20)
    ax.set_xlabel('Variants', labelpad=20)
    ax.set_title('Memory used by pubsub', pad=30, fontsize=20)

    for bar in ax.patches:
        # The text annotation for each bar should be its height.
        bar_value = bar.get_height()
        # Format the text with commas to separate thousands. You can do
        # any type of formatting here though.
        text = f'{bar_value:.2e}'
        # This will give the middle of each bar on the x-axis.
        text_x = bar.get_x() + bar.get_width() / 2
        # get_y() is where the bar starts so we add the height to it.
        text_y = bar.get_y() + bar_value
        # If we want the text to be the same color as the bar, we can
        # get the color like so:
        bar_color = bar.get_facecolor()
        # If you want a consistent color, you can just set it as a constant, e.g. #222222
        ax.text(text_x, text_y, text, ha='center', va='bottom', color=bar_color,
                size=12)

    fig.tight_layout()
    plt.show()
    print()

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

    sns.set_context('talk', font_scale = 0.75)
    fig, ax = plt.subplots(figsize=(12, 8))
    x = np.arange(len(labels))  # the label locations
    width = 0.4  # the width of the bars

    rects1 = ax.bar(x - width/2, mean_values, width, label='mean')
    rects2 = ax.bar(x + width/2, max_values, width, label='max')

    # Add some text for labels, title and custom x-axis tick labels, etc.
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
        # The text annotation for each bar should be its height.
        bar_value = bar.get_height()
        # Format the text with commas to separate thousands. You can do
        # any type of formatting here though.
        text = f'{bar_value:.2f}'
        # This will give the middle of each bar on the x-axis.
        text_x = bar.get_x() + bar.get_width() / 2
        # get_y() is where the bar starts so we add the height to it.
        text_y = bar.get_y() + bar_value
        # If we want the text to be the same color as the bar, we can
        # get the color like so:
        bar_color = bar.get_facecolor()
        # If you want a consistent color, you can just set it as a constant, e.g. #222222
        ax.text(text_x, text_y, text, ha='center', va='bottom', color=bar_color,
                size=12)

    fig.tight_layout()
    plt.show()
    print()

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

    sns.set_context('talk', font_scale = 0.75)
    fig, ax = plt.subplots(figsize=(12, 8))
    x = np.arange(len(labels))  # the label locations
    width = 0.4  # the width of the bars

    rects1 = ax.bar(x - width/2, mean_values, width, label='mean')
    rects2 = ax.bar(x + width/2, max_values, width, label='max')

    # Add some text for labels, title and custom x-axis tick labels, etc.
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
        # The text annotation for each bar should be its height.
        bar_value = bar.get_height()
        # Format the text with commas to separate thousands. You can do
        # any type of formatting here though.
        text = f'{bar_value:.0f}'
        # This will give the middle of each bar on the x-axis.
        text_x = bar.get_x() + bar.get_width() / 2
        # get_y() is where the bar starts so we add the height to it.
        text_y = bar.get_y() + bar_value
        # If we want the text to be the same color as the bar, we can
        # get the color like so:
        bar_color = bar.get_facecolor()
        # If you want a consistent color, you can just set it as a constant, e.g. #222222
        ax.text(text_x, text_y, text, ha='center', va='bottom', color=bar_color,
                size=12)

    fig.tight_layout()
    plt.show()
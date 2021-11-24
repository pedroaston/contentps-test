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

#####################################################################
## plots comparing replication factors and number of subs per user ##
#####################################################################
def final_plot():
    # Collecting metrics
    scout_1_1st = metric_summary("{0}{1}{2}".format("ScoutSubs", 1,"1st"))
    scout_1_2nd = metric_summary("{0}{1}{2}".format("ScoutSubs", 1,"2nd"))
    scout_1_3rd = metric_summary("{0}{1}{2}".format("ScoutSubs", 1,"3rd"))
    scout_1_4th = metric_summary("{0}{1}{2}".format("ScoutSubs", 1,"4th"))

    scout_2_1st = metric_summary("{0}{1}{2}".format("ScoutSubs", 2,"1st"))
    scout_2_2nd = metric_summary("{0}{1}{2}".format("ScoutSubs", 2,"2nd"))
    scout_2_3rd = metric_summary("{0}{1}{2}".format("ScoutSubs", 2,"3rd"))
    scout_2_4th = metric_summary("{0}{1}{2}".format("ScoutSubs", 2,"4th"))

    scout_3_1st = metric_summary("{0}{1}{2}".format("ScoutSubs", 3,"1st"))
    scout_3_2nd = metric_summary("{0}{1}{2}".format("ScoutSubs", 3,"2nd"))
    scout_3_3rd = metric_summary("{0}{1}{2}".format("ScoutSubs", 3,"3rd"))
    scout_3_4th = metric_summary("{0}{1}{2}".format("ScoutSubs", 3,"4th"))

    scout_5_1st = metric_summary("{0}{1}{2}".format("ScoutSubs", 5,"1st"))
    scout_5_2nd = metric_summary("{0}{1}{2}".format("ScoutSubs", 5,"2nd"))
    scout_5_3rd = metric_summary("{0}{1}{2}".format("ScoutSubs", 5,"3rd"))
    scout_5_4th = metric_summary("{0}{1}{2}".format("ScoutSubs", 5,"4th"))

    # Plotting event latency 
    scout_1_event = [scout_1_1st['Event Latency - ScoutSubs11st/mean'], scout_1_2nd['Event Latency - ScoutSubs12nd/mean'],
     scout_1_3rd['Event Latency - ScoutSubs13rd/mean'], scout_1_4th['Event Latency - ScoutSubs14th/mean']]

    scout_2_event = [scout_2_1st['Event Latency - ScoutSubs21st/mean'], scout_2_2nd['Event Latency - ScoutSubs22nd/mean'],
     scout_2_3rd['Event Latency - ScoutSubs23rd/mean'], scout_2_4th['Event Latency - ScoutSubs24th/mean']]
     
    scout_3_event = [scout_3_1st['Event Latency - ScoutSubs31st/mean'], scout_3_2nd['Event Latency - ScoutSubs32nd/mean'],
     scout_3_3rd['Event Latency - ScoutSubs33rd/mean'], scout_3_4th['Event Latency - ScoutSubs34th/mean']]
    
    scout_5_event = [scout_5_1st['Event Latency - ScoutSubs51st/mean'], scout_5_2nd['Event Latency - ScoutSubs52nd/mean'],
     scout_5_3rd['Event Latency - ScoutSubs53rd/mean'], scout_5_4th['Event Latency - ScoutSubs54th/mean']]

    sns.set_context('talk', font_scale = 1)
    fig, ax = plt.subplots(figsize=(12, 8))
    ax.set_title('Avg. Event Latency (at each stage)', pad=30, fontsize=20)
    ax.set_xlabel('# Subscriptions per user', labelpad=20)
    ax.set_ylabel('Time (ms)', labelpad=20)

    t = [1,2,3,4]

    plt.plot(t, scout_1_event, 'rs', label = 'replica-1')
    plt.plot(t, scout_2_event, 'c^', label = 'replica-2')
    plt.plot(t, scout_3_event, 'gs', label = 'replica-3')
    plt.plot(t, scout_5_event, 'k^', label = 'replica-5')

    ax.legend()
    plt.xticks(np.arange(1, 5, 1.0))
    plt.grid(True, alpha = 0.3)
    plt.show()
    print()

    # Plotting event latency 
    scout_1_sub = [scout_1_1st['Sub Latency - ScoutSubs11st/mean'], scout_1_2nd['Sub Latency - ScoutSubs12nd/mean'],
     scout_1_3rd['Sub Latency - ScoutSubs13rd/mean'], scout_1_4th['Sub Latency - ScoutSubs14th/mean']]

    scout_2_sub = [scout_2_1st['Sub Latency - ScoutSubs21st/mean'], scout_2_2nd['Sub Latency - ScoutSubs22nd/mean'],
     scout_2_3rd['Sub Latency - ScoutSubs23rd/mean'], scout_2_4th['Sub Latency - ScoutSubs24th/mean']]
     
    scout_3_sub = [scout_3_1st['Sub Latency - ScoutSubs31st/mean'], scout_3_2nd['Sub Latency - ScoutSubs32nd/mean'],
     scout_3_3rd['Sub Latency - ScoutSubs33rd/mean'], scout_3_4th['Sub Latency - ScoutSubs34th/mean']]
    
    scout_5_sub = [scout_5_1st['Sub Latency - ScoutSubs51st/mean'], scout_5_2nd['Sub Latency - ScoutSubs52nd/mean'],
     scout_5_3rd['Sub Latency - ScoutSubs53rd/mean'], scout_5_4th['Sub Latency - ScoutSubs54th/mean']]

    sns.set_context('talk', font_scale = 1)
    fig, ax = plt.subplots(figsize=(12, 8))
    ax.set_title('Avg. Subscription Latency (at each stage)', pad=30, fontsize=20)
    ax.set_xlabel('# Subscriptions per user', labelpad=20)
    ax.set_ylabel('Time (ms)', labelpad=20)

    plt.plot(t, scout_1_sub, 'rs', label = 'replica-1')
    plt.plot(t, scout_2_sub, 'c^', label = 'replica-2')
    plt.plot(t, scout_3_sub, 'gs', label = 'replica-3')
    plt.plot(t, scout_5_sub, 'k^', label = 'replica-5')

    ax.legend()
    plt.xticks(np.arange(1, 5, 1.0))
    plt.grid(True, alpha = 0.3)
    plt.show()
    print()

    # Plotting memory usage 
    scout_1_memory = [scout_1_1st['Memory used - ScoutSubs11st/mean'],
     scout_1_1st['Memory used - ScoutSubs11st/mean']+scout_1_2nd['Memory used - ScoutSubs12nd/mean'],
     scout_1_1st['Memory used - ScoutSubs11st/mean']+scout_1_2nd['Memory used - ScoutSubs12nd/mean']+scout_1_3rd['Memory used - ScoutSubs13rd/mean'],
     scout_1_1st['Memory used - ScoutSubs11st/mean']+scout_1_2nd['Memory used - ScoutSubs12nd/mean']+scout_1_3rd['Memory used - ScoutSubs13rd/mean']+scout_1_4th['Memory used - ScoutSubs14th/mean']]
     
    scout_2_memory = [scout_2_1st['Memory used - ScoutSubs21st/mean'],
     scout_2_1st['Memory used - ScoutSubs21st/mean']+scout_2_2nd['Memory used - ScoutSubs22nd/mean'],
     scout_2_1st['Memory used - ScoutSubs21st/mean']+scout_2_2nd['Memory used - ScoutSubs22nd/mean']+scout_2_3rd['Memory used - ScoutSubs23rd/mean'],
     scout_2_1st['Memory used - ScoutSubs21st/mean']+scout_2_2nd['Memory used - ScoutSubs22nd/mean']+scout_2_3rd['Memory used - ScoutSubs23rd/mean']+scout_2_4th['Memory used - ScoutSubs24th/mean']]
    
    scout_3_memory = [scout_3_1st['Memory used - ScoutSubs31st/mean'],
     scout_3_1st['Memory used - ScoutSubs31st/mean']+scout_3_2nd['Memory used - ScoutSubs32nd/mean'],
     scout_3_1st['Memory used - ScoutSubs31st/mean']+scout_3_2nd['Memory used - ScoutSubs32nd/mean']+scout_3_3rd['Memory used - ScoutSubs33rd/mean'],
     scout_3_1st['Memory used - ScoutSubs31st/mean']+scout_3_2nd['Memory used - ScoutSubs32nd/mean']+scout_3_3rd['Memory used - ScoutSubs33rd/mean']+scout_3_4th['Memory used - ScoutSubs34th/mean']]

    scout_5_memory = [scout_5_1st['Memory used - ScoutSubs51st/mean'],
     scout_5_2nd['Memory used - ScoutSubs52nd/mean']+scout_5_1st['Memory used - ScoutSubs51st/mean'],
     scout_5_2nd['Memory used - ScoutSubs52nd/mean']+scout_5_1st['Memory used - ScoutSubs51st/mean']+scout_5_3rd['Memory used - ScoutSubs53rd/mean'],
     scout_5_2nd['Memory used - ScoutSubs52nd/mean']+scout_5_1st['Memory used - ScoutSubs51st/mean']+scout_5_3rd['Memory used - ScoutSubs53rd/mean']+scout_5_4th['Memory used - ScoutSubs54th/mean']]

    sns.set_context('talk', font_scale = 1)
    fig, ax = plt.subplots(figsize=(12, 8))
    ax.set_title('Avg. Memory Used (cumulatively)', pad=30, fontsize=20)
    ax.set_xlabel('# Subscriptions per user', labelpad=20)
    ax.set_ylabel('# Memory (MB)', labelpad=20)

    plt.plot(t, scout_1_memory, 'r--', label = 'replica-1')
    plt.plot(t, scout_2_memory, 'c', label = 'replica-2')
    plt.plot(t, scout_3_memory, 'g--', label = 'replica-3')
    plt.plot(t, scout_5_memory, 'k', label = 'replica-5')

    ax.legend()
    plt.xticks(np.arange(1, 5, 1.0))
    plt.grid(True, alpha = 0.3)
    plt.show()
    print()

    # Plotting cpu usage 
    scout_1_cpu = [scout_1_1st['CPU used - ScoutSubs11st/mean'],
     scout_1_2nd['CPU used - ScoutSubs12nd/mean']+scout_1_1st['CPU used - ScoutSubs11st/mean'],
     scout_1_3rd['CPU used - ScoutSubs13rd/mean']+scout_1_2nd['CPU used - ScoutSubs12nd/mean']+scout_1_1st['CPU used - ScoutSubs11st/mean'],
     scout_1_4th['CPU used - ScoutSubs14th/mean']+scout_1_3rd['CPU used - ScoutSubs13rd/mean']+scout_1_2nd['CPU used - ScoutSubs12nd/mean']+scout_1_1st['CPU used - ScoutSubs11st/mean']]
     
    scout_2_cpu = [scout_2_1st['CPU used - ScoutSubs21st/mean'], 
     scout_2_2nd['CPU used - ScoutSubs22nd/mean']+scout_2_1st['CPU used - ScoutSubs21st/mean'],
     scout_2_3rd['CPU used - ScoutSubs23rd/mean']+scout_2_2nd['CPU used - ScoutSubs22nd/mean']+scout_2_1st['CPU used - ScoutSubs21st/mean'],
     scout_2_4th['CPU used - ScoutSubs24th/mean']+scout_2_3rd['CPU used - ScoutSubs23rd/mean']+scout_2_2nd['CPU used - ScoutSubs22nd/mean']+scout_2_1st['CPU used - ScoutSubs21st/mean']]
    
    scout_3_cpu = [scout_3_1st['CPU used - ScoutSubs31st/mean'],
     scout_3_2nd['CPU used - ScoutSubs32nd/mean']+scout_3_1st['CPU used - ScoutSubs31st/mean'],
     scout_3_3rd['CPU used - ScoutSubs33rd/mean']+scout_3_2nd['CPU used - ScoutSubs32nd/mean']+scout_3_1st['CPU used - ScoutSubs31st/mean'],
     scout_3_4th['CPU used - ScoutSubs34th/mean']+scout_3_3rd['CPU used - ScoutSubs33rd/mean']+scout_3_2nd['CPU used - ScoutSubs32nd/mean']+scout_3_1st['CPU used - ScoutSubs31st/mean']]

    scout_5_cpu = [scout_5_1st['CPU used - ScoutSubs51st/mean'],
     scout_5_1st['CPU used - ScoutSubs51st/mean']+scout_5_2nd['CPU used - ScoutSubs52nd/mean'],
     scout_5_3rd['CPU used - ScoutSubs53rd/mean']+scout_5_1st['CPU used - ScoutSubs51st/mean']+scout_5_2nd['CPU used - ScoutSubs52nd/mean'],
     scout_5_4th['CPU used - ScoutSubs54th/mean']+scout_5_3rd['CPU used - ScoutSubs53rd/mean']+scout_5_1st['CPU used - ScoutSubs51st/mean']+scout_5_2nd['CPU used - ScoutSubs52nd/mean']]

    sns.set_context('talk', font_scale = 1)
    fig, ax = plt.subplots(figsize=(12, 8))
    ax.set_title('Avg. CPU User-time used (cumulatively)', pad=30, fontsize=20)
    ax.set_xlabel('# Subscriptions per user', labelpad=20)
    ax.set_ylabel('Time (s)', labelpad=20)

    plt.plot(t, scout_1_cpu, 'r--', label = 'replica-1')
    plt.plot(t, scout_2_cpu, 'c', label = 'replica-2')
    plt.plot(t, scout_3_cpu, 'g--', label = 'replica-3')
    plt.plot(t, scout_5_cpu, 'k', label = 'replica-5')

    ax.legend()
    plt.xticks(np.arange(1, 5, 1.0))
    plt.grid(True, alpha = 0.3)
    plt.show()

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
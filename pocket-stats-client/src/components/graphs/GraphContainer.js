import React, {Component} from 'react';
import moment from 'moment';
import Graph from './Graph.js';
import GraphMenu from './GraphMenu.js';
import { GraphTypes, ChartType } from '../../constants/graphTypes.js';
import { LineGraphData } from '../../constants/graphData';
import { GraphColours } from '../../constants/colours';

class GraphContainer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      currentGraph: GraphTypes.ARTICLES_READ.name,
      currentChartType: ChartType.LINE
    }
  }

  onMenuItemClicked = (graphType) => {
    this.setState({
      currentGraph: graphType.name,
      currentChartType: graphType.type,
    });
  }

  getGraphData() {
    const {itemisedDate} = this.props;

    let graphData = [];

    // Create the date labels
    const labels = Object.keys(itemisedDate).map((key) => (
      moment.unix(key).format("D/MMM")
    ));

    console.log(GraphTypes.ARTICLES_ADDED.name)
    switch (this.state.currentGraph) {
      case GraphTypes.ARTICLES_READ.name:
        graphData = Object.keys(itemisedDate).map((key) => (
          itemisedDate[key].articles_read
        ))
        break;
      case GraphTypes.ARTICLES_ADDED.name:
        graphData = Object.keys(itemisedDate).map((key) => (
          itemisedDate[key].articles_added
        ))
        break;
      case GraphTypes.WORDS_READ.name:
        graphData = Object.keys(itemisedDate).map((key) => (
          itemisedDate[key].words_read
        ))
        break;
      case GraphTypes.WORDS_ADDED.name:
        graphData = Object.keys(itemisedDate).map((key) => (
          itemisedDate[key].words_added
        ))
        break;
      default:
        console.log("Nani")
        break;
    }

    return {
      labels: labels,
      data: graphData
    };
  }

  // getAndSortTags gets the relevant stat for the current graph type and sorts
  // the data from largest to smallest
  getAndSortTags(itemisedTags) {
    let tags = [];
    switch (this.state.currentGraph) {
      case GraphTypes.TAGS_READ.name:
      tags = Object.keys(itemisedTags).map((key) => ({
        name: key,
        value: itemisedTags[key].articles_read
      }));
      break;

      case GraphTypes.TAGS_TIME.name:
      console.log(itemisedTags)
      tags = Object.keys(itemisedTags).map((key) => ({
        name: key,
        value: itemisedTags[key].time_reading
      }));
      break;

      default:
      console.error("How did we get an unknown graph type?")
      return;
    }

    return tags.sort((a, b) => b.value - a.value);
  }

  // getDonutGraphData creates the data structure needed to create a donut (pie)
  // chart
  getDonutGraphData() {
    const { itemisedTags } = this.props;

    // We need to sort the tags, as it looks better
    const tagArray = this.getAndSortTags(itemisedTags);

    // Get the names and values from the sorted array
    const labels = tagArray.map(t => t.name);
    const graphData = tagArray.map(t => t.value);

    // Loop through the 5 colours we have and assign them to a tag
    let colours = [];
    const noTags = labels.length;
    for (let i=0; i < noTags; i++) {
      colours.push(GraphColours[i%15]);
    }

    const data = {
    	labels: labels,
    	datasets: [{
    		data: graphData,
    		backgroundColor: colours,
    		hoverBackgroundColor: colours
    	}]
    };

    return data;
  }

  render() {
    let data;
    if (this.state.currentChartType === ChartType.PIE) {
      data = this.getDonutGraphData();
    } else {
      const d3Data = this.getGraphData();
      data = LineGraphData;
      data.labels = d3Data.labels;
      data.datasets[0].label = this.state.currentGraph;
      data.datasets[0].data = d3Data.data;
      console.log(d3Data);
    }

    return(
      <div>
        <div className="row">
          <div className="col-lg-2">
            <GraphMenu onClick={this.onMenuItemClicked} />
          </div>
          <div className='col-lg-10'>
            <Graph data={data} graphType={this.state.currentChartType}/>
          </div>
        </div>
      </div>
    );
  }
}

export default GraphContainer;

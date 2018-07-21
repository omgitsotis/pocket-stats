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

  getDonutGraphData() {
    const {itemisedTags} = this.props;

    let graphData = [];
    let colours = [];

    const labels = Object.keys(itemisedTags);
    const noTags = labels.length;

    // Loop through the 5 colours we have and assign them to a tag
    for (let i=0; i < noTags; i++) {
      colours.push(GraphColours[i%5]);
    }

    graphData = Object.keys(itemisedTags).map((key) => (
      itemisedTags[key].articles_read
    ))

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
    if (this.state.currentGraph === GraphTypes.TAGS_READ.name) {
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

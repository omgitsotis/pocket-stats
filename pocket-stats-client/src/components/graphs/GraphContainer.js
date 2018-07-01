import React, {Component} from 'react';
import moment from 'moment';
import Graph from './Graph.js';
import GraphMenu from './GraphMenu.js'
import GraphTypes from '../../constants/graphTypes.js'

class GraphContainer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      currentGraph: GraphTypes.ARTICLES_READ
    }
  }

  onMenuItemClicked = (graphType) => {
    this.setState({currentGraph: graphType})
  }

  getGraphData() {
    const {totals, itemised} = this.props;

    let graphData = [];

    // Create the date labels
    const labels = Object.keys(itemised).map((key) => (
      moment.unix(key).format("D/MMM")
    ));

    switch (this.state.currentGraph) {
      case GraphTypes.ARTICLES_READ:
        graphData = Object.keys(itemised).map((key) => (
          itemised[key].articles_read
        ))
        break;
      case GraphTypes.ARTICLES_ADDED:
        graphData = Object.keys(itemised).map((key) => (
          itemised[key].articles_added
        ))
        break;
      case GraphTypes.WORDS_READ:
        graphData = Object.keys(itemised).map((key) => (
          itemised[key].words_read
        ))
        break;
      case GraphTypes.ARTICLES_ADDED:
        graphData = Object.keys(itemised).map((key) => (
          itemised[key].words_added
        ))
        break;
      default:
        break;
    }

    return {
      labels: labels,
      data: graphData
    };
  }

  render() {
    const d3Data = this.getGraphData()

    // TODO: Move this data somewhere
    const data = {
      labels: d3Data.labels,
      datasets: [{
        label: this.state.currentGraph,
        fill: false,
        lineTension: 0.1,
        backgroundColor: 'rgba(75,192,192,0.4)',
        borderColor: 'rgba(75,192,192,1)',
        borderCapStyle: 'butt',
        borderDash: [],
        borderDashOffset: 0.0,
        borderJoinStyle: 'miter',
        pointBorderColor: 'rgba(75,192,192,1)',
        pointBackgroundColor: '#fff',
        pointBorderWidth: 1,
        pointHoverRadius: 5,
        pointHoverBackgroundColor: 'rgba(75,192,192,1)',
        pointHoverBorderColor: 'rgba(220,220,220,1)',
        pointHoverBorderWidth: 2,
        pointRadius: 1,
        pointHitRadius: 10,
        data: d3Data.data
      }]
    };

    return(
      <div>
        <div className="row">
          <div className="col-lg-2">
            <GraphMenu onClick={this.onMenuItemClicked} />
          </div>
          <div className='col-lg-10'>
            <Graph data={data} />
          </div>
        </div>
      </div>
    );
  }
}

export default GraphContainer;

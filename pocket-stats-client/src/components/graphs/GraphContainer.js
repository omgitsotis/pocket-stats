import React, {Component} from 'react';
import moment from 'moment';
import Graph from './Graph.js';

class GraphContainer extends Component {
  render() {
    const totals = this.props.totals;
    const itemised = this.props.itemised;

    let labels = [];
    let atsRead = [];

    Object.keys(itemised).forEach(function(key) {
      labels.push(moment.unix(key).format("D/MMM"));
      atsRead.push(itemised[key].articles_read);
    });

    const data = {
      labels: labels,
      datasets: [{
        label: 'Articles Read',
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
        data: atsRead
      }]
    };

    return(
      <div>
        <div className="row">
          <div className="col-lg-2"><p>Menu coming soon</p></div>
          <div className='col-lg-10'>
            <Graph data={data} />
          </div>
        </div>
      </div>
    );
  }
}

export default GraphContainer;

import React from 'react'
import { Pie, Line } from 'react-chartjs-2';
import { ChartType } from '../../constants/graphTypes.js';

function Graph(props) {
  let graph;
  switch (props.graphType) {
    case ChartType.LINE:
      graph = <Line data={props.data} height={100}/>
      break;
    case ChartType.PIE:
      graph = <Pie data={props.data} height={200}/>
      break;
    default:
      console.log("What the fuck Otis?");
      break;
  }

  return (
    <div>
      <div className="row">
        <div className='col-lg'>
          {graph}
        </div>
      </div>
    </div>
  )
}

export default Graph

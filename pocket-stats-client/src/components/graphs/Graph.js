import React from 'react'
import {Line} from 'react-chartjs-2';

function Graph(props) {
  return (
    <div>
      <div className="row">
        <div className='col-lg'>
          <Line data={props.data} height={100}/>
        </div>
      </div>
    </div>
  )
}

export default Graph

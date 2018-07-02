import React from 'react'
import GraphTypes from '../../constants/graphTypes.js'

function GraphMenuItem({graphType, onClick}) {
  return (
    <li onClick={() => onClick(graphType)}>{graphType}</li>
  )
}

function GraphMenu({onClick}) {
  return (
    <ul>
      {Object.keys(GraphTypes).map((key, i) => (
        <GraphMenuItem key={i} onClick={onClick} graphType={GraphTypes[key]} />
      ))}
    </ul>
  )
}



export default GraphMenu

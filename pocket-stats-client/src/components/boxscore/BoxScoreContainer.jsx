import React, {Component} from 'react';
import BoxScore from './BoxScore.jsx';

class BoxScoreContainer extends Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <BoxScore totals={this.props.totals}/>
        )
    }
}

export default BoxScoreContainer;

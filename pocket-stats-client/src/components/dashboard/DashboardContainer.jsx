import React, {Component} from 'react';
import Dashboard from './Dashboard.jsx';

class DashboardContainer extends 'Component'{
    render() {
        return (
            <Dashboard
                totals={this.props.totals}
                itemised={this.props.itemised}
                lastUpdated={this.props.lastUpdated}
            />
        )
    }
}

export default DashboardContainer;

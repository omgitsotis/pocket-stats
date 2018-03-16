import React, {Component} from 'react';
import Dashboard from './Dashboard.jsx';

class DashboardContainer extends 'Component'{
    render() {
        return (
            <Dashboard {...this.props} />
        )
    }
}

export default DashboardContainer;

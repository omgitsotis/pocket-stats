import React, {Component} from 'react';
import './dashboard.css';

class Dashboard extends Component {
    render() {
        return(
            <div className='row'>
                <div className='col-lg-2'>
                    <p>Sidebar</p>
                </div>
                <div className='col-lg-10'>
                    <div className='row stat-grid'>
                        <div className="col-lg-5">
                            <div className='stats-section'>
                                <div className="dummy"></div>
                                <p className='thumbnail purple'>Homepage</p>
                            </div>
                        </div>
                        <div className="col-lg-5">
                            <div className='stats-section'>
                                <div className="dummy"></div>
                                <p className='thumbnail purple'>Homepage</p>
                            </div>
                        </div>
                        <div className="col-lg-5">
                            <div className='stats-section'>
                                <div className="dummy"></div>
                                <p className='thumbnail purple'>Homepage</p>
                            </div>
                        </div>
                        <div className="col-lg-5">
                            <div className='stats-section'>
                                <div className="dummy"></div>
                                <p className='thumbnail purple'>Homepage</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

export default Dashboard;

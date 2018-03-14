import React, {Component} from 'react';
import './dashboard.css';
import {Line} from 'react-chartjs-2';

class Dashboard extends Component {
    render() {
        const data = {
            labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
            datasets: [{
                label: 'My First dataset',
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
                data: [65, 59, 80, 81, 56, 55, 40]
            }]
        };

        return(
            <div>
                <div className='row'>
                    <div className='col-lg'>
                        <nav class="navbar navbar-light bg-light">
                            <span class="navbar-brand mb-0 h1">Navbar</span>
                        </nav>
                    </div>
                </div>
                <div className='row'>
                    <div className="col-lg">
                        <div className='card-deck'>
                            <div className='card text-center stat-box'>
                                <h5 className='card-title'>Articles Add</h5>
                                <p className='card-text'>30</p>
                            </div>
                            <div className='card text-center stat-box'>
                                <h5 className='card-title'>Articles Read</h5>
                                <p className='card-text'>15</p>
                            </div>
                            <div className='card text-center stat-box'>
                                <h5 className='card-title'>Words add</h5>
                                <p className='card-text'>30</p>
                            </div>
                            <div className='card text-center stat-box'>
                                <h5 className='card-title'>Words Read</h5>
                                <p className='card-text'>30</p>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="row">
                    <div className='col-lg'>
                        <Line data={data} height={100}/>
                    </div>
                </div>
            </div>
        );
    }
}

export default Dashboard;

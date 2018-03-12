import React, {Component} from 'react';
import './menu.css'

class Menu extends Component {
    render() {
        return(
            <div className='col-lg-6 offset-lg-3'>
                <h1>Main Menu</h1>
                <div className='btn-row'>
                    <button className='btn btn-primary btn-block'
                        onClick={this.props.onInitClick}>Initalise</button>
                </div>
                <div className='btn-row'>
                    <button className='btn btn-primary btn-block'
                        onClick={ () => this.onClick('overview') }>Overview</button>
                </div>
                <div className='btn-row'>
                    <button className='btn btn-primary btn-block'
                        onClick={ () => this.onClick('update') }>Update</button>
                </div>
            </div>
        )
    }
}

export default Menu;

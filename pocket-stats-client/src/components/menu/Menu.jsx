import React, {Component} from 'react';
import InitSection from './../initialisation/init.jsx';

class Menu extends Component {
    constructor(props) {
        super(props);
        this.state = {
            currentPage: 'menu'
        };
    }

    onClick = (page) => {
        console.log(page);
        this.setState({currentPage: page});
    }

    render() {
        let section;
        switch (this.state.currentPage) {
            case 'init':
                section =
                    <InitSection
                        onBackClick={ () => this.onClick('menu')}
                        {...this.props}
                    />;
                break;
            default:
                section =
                    <div>
                        <button onClick={ () => this.onClick('init') }>Initalise</button>
                        <button onClick={ () => this.onClick('overview') }>Overview</button>
                    </div>;
                break;
        }
        return (
            <div>
                {section}
            </div>
        )
    }
}

export default Menu;

import React, {Component} from 'react';
import {Helmet} from "react-helmet";
import './App.css';
import ButtonPad from "./ButtonPad";
import {BrowserRouter, Link, Route, Switch} from 'react-router-dom'
import Virtual from "./Virtual";
import ButtonEdit from "./button-edit/ButtonEdit";

// Home page is just a bunch of links for now
const Home = () => (
    <nav>
        <Link to="/buttons">Buttons</Link>
        <hr/>
        <Link to="/virtual">Virtual</Link>
    </nav>
);

// The main application navigates between components TODO trap button long press and move to button edit...
class App extends Component {

    render() {
        return (
            <div className="App">
                <Helmet>
                    <meta name="apple-mobile-web-app-capable" content="yes"/>
                    <title>Bright Light</title>
                </Helmet>
                <BrowserRouter>
                    <Switch>
                        <Route path="/buttons" component={ButtonPad}/>
                        <Route path="/button-edit" component={ButtonEdit}/>
                        <Route path="/virtual" component={Virtual}/>
                        <Route path="/" component={Home}/>
                    </Switch>
                </BrowserRouter>
            </div>
        );
    }
}

export default App;

import React, {Component} from 'react';
import './App.css';
import ButtonPad from "./buttons/ButtonPad";
import {BrowserRouter, Link, Route, Routes, useNavigate, useLocation} from 'react-router-dom'
import ButtonEditor from "./button-edit/ButtonEditor";
import {getStaticData} from "./server-proxy/staticData";
import {getButtons, saveButtons} from "./server-proxy/buttons";
import {runAnimation} from "./server-proxy/animation";
import Virtual from "./virtual/Virtual";
import {OpenWebSocket} from "./server-proxy/webSocket";

// Home page is just a bunch of links for now
const Home = () => (
    <nav>
        <Link to="/buttons">Buttons</Link>
        <hr/>
        <Link to="/virtual">Virtual</Link>
        <hr/>
    </nav>
);

// Wrapper providing useNavigate to ButtonPad (class component can't use hooks)
function ButtonPadRoute({allButtons, activeButtonKey, onButtonTap, onButtonPress}) {
    const navigate = useNavigate();
    return <ButtonPad allButtons={allButtons}
                      activeButtonKey={activeButtonKey}
                      onButtonTap={onButtonTap}
                      onButtonPress={key => onButtonPress(navigate, key)}/>;
}

// Wrapper providing useNavigate/useLocation to ButtonEditor (class component can't use hooks)
function ButtonEditorRoute({allButtons, allAnimations, allSegments, userSegments, onButtonChanged, onOk}) {
    const navigate = useNavigate();
    const location = useLocation();
    const history = {location, goBack: () => navigate(-1)};
    return <ButtonEditor allButtons={allButtons}
                         allAnimations={allAnimations}
                         allSegments={allSegments}
                         userSegments={userSegments}
                         onButtonChanged={onButtonChanged}
                         onOk={onOk}
                         history={history}/>;
}

// The main application navigates between components
export default class App extends Component {

    constructor(props) {
        super(props);
        this.state =
            {
                userButtons: undefined,
                userSegments: undefined,
                allSegments: undefined,
                allAnimations: undefined,
                activeButtonKey: 0,
                buttonPadVersion: 1
            };
    }

    // Get static data from server when mounting
    componentDidMount() {

        console.info('getting static data');
        getStaticData(
            sd => this.setState({allSegments: sd.segments, allAnimations: sd.animations}),
            (xhr) => console.error(xhr));

        console.info('getting button configuration');
        this.getButtonConfig();

        // Subscribe to button state updates (we get immediate update on connection)
        this.ws = OpenWebSocket('ButtonState', bs => {
            console.info('button state changed to ', bs);

            if (this.state.activeButtonKey !== bs.ActiveButtonKey) {
                console.info('updating active button');
                this.setState({activeButtonKey: bs.ActiveButtonKey});
            }

            if (this.state.buttonPadVersion !== bs.ButtonPadVersion) {
                console.info('updating button configuration');
                this.setState({buttonPadVersion: bs.ButtonPadVersion});
                this.getButtonConfig();
            }
        });
    }

    getButtonConfig() {
        getButtons(
            cfg => this.setState({userSegments: cfg.segments, userButtons: cfg.buttons}),
            (xhr) => console.error(xhr));
    }

    // Close web socket
    componentWillUnmount() {
        if (this.ws !== undefined)
            this.ws.close()
    }

    onButtonChanged(button) {
        let userButtons = this.state.userButtons.map(b => b.key === button.key ? button : b);
        this.setState((props, state) => {
            return {...state, userButtons: userButtons}
        });
        runAnimation(button);
    }

    onSaveButtonEdit() {
        console.info('onSaveButtonEdit');
        saveButtons(
            {segments: this.state.userSegments, buttons: this.state.userButtons},
            () => console.info('button state saved'),
            (xhr) => console.error(xhr)); // TODO: Report errors to user
    }

    findButton(key) {
        return this.state.userButtons.find(x => x.key === key);
    }

    onButtonTap(key) {
        runAnimation(this.findButton(key));
    }

    onButtonPress(navigate, key) {
        runAnimation(this.findButton(key));
        navigate('/button-edit', {state: {'buttonKey': key}});
    }

    render() {
        return (
            <div className="App" onContextMenu={e => e.preventDefault()}>
                <BrowserRouter>
                    <Routes>
                        <Route path="/buttons" element={
                            <ButtonPadRoute
                                allButtons={this.state.userButtons}
                                activeButtonKey={this.state.activeButtonKey}
                                onButtonTap={key => this.onButtonTap(key)}
                                onButtonPress={(navigate, key) => this.onButtonPress(navigate, key)}
                            />
                        }/>
                        <Route path="/button-edit" element={
                            <ButtonEditorRoute
                                allButtons={this.state.userButtons}
                                allAnimations={this.state.allAnimations}
                                allSegments={this.state.allSegments}
                                userSegments={this.state.userSegments}
                                onButtonChanged={button => this.onButtonChanged(button)}
                                onOk={() => this.onSaveButtonEdit()}
                            />
                        }/>
                        <Route path="/virtual" element={<Virtual/>}/>
                        <Route path="/" element={<Home/>}/>
                    </Routes>
                </BrowserRouter>
            </div>
        );
    }
}

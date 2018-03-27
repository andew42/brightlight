import * as React from "react";
import {Fragment} from "react";
import NameEditor from "./NameEditor";
import './ButtonEdit.css';
import {saveButtons} from "../server-proxy/buttons";
import {Button} from "semantic-ui-react";
import LedSegment from "./LedSegment";
import LedSegmentList from "./LedSegmentList";

// Shows a list of editors that change to suit the animation
export default class ButtonEdit extends React.Component {
    constructor(props) {
        super(props);

        // We expect to be passed all buttons and button to edit in location state
        this.buttons = props.history.location.state.buttons;
        this.button = props.history.location.state.button;

        // Copy the initial state of the button we are editing
        let buttonSegmentsCopy = JSON.parse(JSON.stringify(this.button.segments));
        this.state = {
            name: this.button.name,
            segments: buttonSegmentsCopy,
            selectedSegments: buttonSegmentsCopy.map(s => s.segment),
            isOpen: false
        };

        // Callback bindings for this
        this.onOK = this.onOK.bind(this);
        this.toggleSelectedSegment = this.toggleSelectedSegment.bind(this);
        this.editSegmentListCancel = this.editSegmentListCancel.bind(this);
        this.editSegmentListOk = this.editSegmentListOk.bind(this);
    }

    onOK() {
        // Update button object and save to server
        this.button.name = this.state.name;
        // TODO: Better error handling
        saveButtons(this.buttons, () => this.props.history.goBack(), e => console.error(e));
    }

    onRemoveSegment(segment) {
        console.info('removing segment ' + segment.segment);
        this.setState({segments: this.state.segments.filter(seg => seg !== segment)})
    }

    toggleSelectedSegment(segName) {
        console.info('Toggle:' + segName);
        if (this.state.selectedSegments.includes(segName))
            this.setState({selectedSegments: this.state.selectedSegments.filter(n => n !== segName)});
        else {
            let selectedSegmentCopy = JSON.parse(JSON.stringify(this.state.selectedSegments));
            selectedSegmentCopy.push(segName);
            this.setState({selectedSegments: selectedSegmentCopy});
        }
    }

    editSegmentListCancel() {
        this.setState({
            selectedSegments: this.state.segments.map(s => s.segment),
            isOpen: false
        });
    }

    editSegmentListOk() {
        // TODO
        this.setState({isOpen: false});
    }

    render() {
        //let selectedSegments = this.state.segments.map(s => s.segment);
        let key = 1;
        return <div className="button-edit-editor-list">
            <Fragment>
                <NameEditor name={this.state.name} onNameChanged={name => this.setState({name: name})}/>

                {this.state.segments.map(segment => (
                    <LedSegment key={key++} segment={segment} onRemove={seg => this.onRemoveSegment(seg)}/>))}

                <div className='ok-cancel-container'>
                    <LedSegmentList isOpen={this.state.isOpen}
                                    onCancel={this.editSegmentListCancel}
                                    onOk={this.editSegmentListOk}
                                    selectedItems={this.state.selectedSegments}
                                    toggleSelectedSegment={seg => this.toggleSelectedSegment(seg)}
                                    trigger={<Button icon='plus' circular floated='left'
                                                     onClick={() => this.setState({isOpen: true})}/>}/>
                    <Button primary onClick={this.onOK}>OK</Button>
                    <Button secondary onClick={() => this.props.history.goBack()}>Cancel</Button>
                </div>
            </Fragment>
        </div>
    }
}

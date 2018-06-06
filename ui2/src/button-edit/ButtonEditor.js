import * as React from "react";
import {Fragment} from "react";
import './ButtonEditor.css';
import {Button} from "semantic-ui-react";
import {LedSegmentEditor} from "./LedSegmentEditor";
import {NameEditor} from "./NameEditor";
import LedSegmentChooser from "./LedSegmentChooser";

// Shows a list of editors that change to suit the animation
export default class ButtonEditor extends React.Component {

    componentDidMount() {
        // Make a copy of the initial button state in case we cancel
        let allButtons = this.props.allButtons;
        let buttonKey = this.props.history.location.state.buttonKey;
        let button = allButtons.find(x => x.key === buttonKey);
        if (button !== undefined)
            this.initialButton = JSON.parse(JSON.stringify(button));
    }

    toggleSelectedSegment(seg) {
        let allButtons = this.props.allButtons;
        let buttonKey = this.props.history.location.state.buttonKey;
        let button = allButtons.find(x => x.key === buttonKey);
        // Make a note of initial button in case we cancel
        if (this.initialButton === undefined)
            this.initialButton = {...button};
        if (button.segments.find(s => s.name === seg.name) === undefined) {
            // Add new segment to the button
            let newSegments = button.segments.map(x => x);
            newSegments.push({...seg, animation: "Static", params: "#3f3f3f"});
            newSegments.sort((a, b) => a.z - b.z);
            this.props.onButtonChanged({
                ...button,
                segments: newSegments
            });
        }
        else {
            // Remove segment from button
            this.props.onButtonChanged({
                ...button,
                segments: button.segments.filter(s => s.name !== seg.name)
            });
        }
    }

    render() {
        let allButtons = this.props.allButtons;
        let allAnimations = this.props.allAnimations;
        let allSegments = this.props.allSegments;
        let buttonKey = this.props.history.location.state.buttonKey;
        let button = allButtons.find(x => x.key === buttonKey);
        let animationNames = allAnimations.map(n => ({'text': n.name, 'value': n.name}));
        let otherButtonNames = allButtons.filter(b => b.key !== buttonKey).map(b => b.name);
        let key = 1;
        return <div className="button-editor-editor-list">
            <Fragment>
                <NameEditor name={button.name}
                            onNameChanged={newName => this.props.onButtonChanged({
                                ...button,
                                name: newName
                            })}
                            error={otherButtonNames.find(x => x.toUpperCase() === button.name.toUpperCase()) ?
                                'Name already exists' :
                                undefined}/>

                {button.segments.map(segment => (
                    <LedSegmentEditor key={key++}
                                      segment={segment}
                                      allAnimationNames={animationNames}
                                      onRemove={seg => this.props.onButtonChanged({
                                          ...button,
                                          segments: button.segments.filter(s => s.name !== seg.name)
                                      })}
                                      onSegmentChanged={seg => this.props.onButtonChanged({
                                          ...button,
                                          segments: button.segments.map(s => s.name === seg.name ? seg : s)
                                      })}/>))}

                <div className='button-editor-ok-cancel-container'>
                    <LedSegmentChooser allSegments={allSegments}
                                       checkedSegmentNames={button.segments.map(s => s.name)}
                                       onOk={() => this.initialButton = undefined}
                                       onCancel={() => {
                                           this.props.onButtonChanged(this.initialButton);
                                           this.initialButton = undefined;
                                       }}
                                       toggleCheckedSegment={seg => this.toggleSelectedSegment(seg)}
                                       trigger={<Button icon='plus'
                                                        circular
                                                        floated='left'/>}/>
                    <Button primary onClick={() => {
                        this.props.onOk();
                        this.props.history.goBack();
                    }} content='OK'/>
                    <Button secondary onClick={() => {
                        this.props.onButtonChanged(this.initialButton);
                        this.props.history.goBack();
                    }} content='Cancel'/>
                </div>
            </Fragment>
        </div>
    }
}

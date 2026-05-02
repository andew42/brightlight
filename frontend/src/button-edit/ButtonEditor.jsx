import * as React from "react";
import {Fragment} from "react";
import './ButtonEditor.css';
import {LedSegmentEditor} from "./LedSegmentEditor";
import {NameEditor} from "./NameEditor";
import LedSegmentChooser from "./LedSegmentChooser";
import {Button} from "semantic-ui-react";

// Shows a list of editors that change to suit the animation
export default class ButtonEditor extends React.Component {

    constructor(props) {
        super(props);
        this.state = {selectedSegments: []};
    };

    componentDidMount() {
        // Make a copy of the initial button state in case we cancel
        let allButtons = this.props.allButtons;
        let buttonKey = this.props.history.location.state.buttonKey;
        let button = allButtons.find(x => x.key === buttonKey);
        if (button !== undefined)
            this.preEditButton = {...button};
    }

    toggleSelectedSegment(seg) {
        let allButtons = this.props.allButtons;
        let buttonKey = this.props.history.location.state.buttonKey;
        let button = allButtons.find(x => x.key === buttonKey);

        // Make a note of initial button in case we cancel
        if (this.preSegChooserButton === undefined)
            this.preSegChooserButton = {...button};

        // Work out top z order for this segment name as there
        // may be zero or more segments with this name
        let topz = button.segments.filter(s => s.name === seg.name).map(s => s.z).reduce((_, current, index, array) => Math.max(current, array[index]), undefined);

        // Is this a selection or a removal (check collection for this add popup)
        if (this.state.selectedSegments.find(s => s === seg.name) === undefined) {

            // Add NEW segment to check collection and button
            let newSelectedSegments = this.state.selectedSegments.map(x => x);
            newSelectedSegments.push(seg.name);
            this.setState({selectedSegments: newSelectedSegments});

            // Make a copy of the new segments
            let newSegments = button.segments.map(x => x);

            // New segment to add with a static animation
            let newSegment = {
                ...seg, animation: "Static", params: [
                    {"key": 20, "type": "colour", "label": "Colour", "value": {"r": 31, "g": 31, "b": 31}}]
            };

            // Adjust new segment Z order if necessary (already segment(s) with this name)
            if (topz !== undefined)
                newSegment.z = topz + 0.1;

            // Add new segment to the new segments copy
            newSegments.push(newSegment);

            // Sort the segments into z order
            newSegments.sort((a, b) => a.z - b.z);

            // Update the button so we see changes immediately
            this.props.onButtonChanged({
                ...button,
                segments: newSegments
            });
        } else {
            // REMOVE segment from check collection...
            this.setState({selectedSegments: this.state.selectedSegments.filter(s => s !== seg.name)});

            // ...and button
            this.props.onButtonChanged({
                ...button,
                segments: button.segments.filter(s => s.name !== seg.name || s.z !== topz)
            });
        }
    }

    render() {
        let allButtons = this.props.allButtons;
        let allAnimations = this.props.allAnimations;
        let allSegments = this.props.allSegments;
        let userSegments = this.props.userSegments;
        let buttonKey = this.props.history.location.state.buttonKey;
        let button = allButtons.find(x => x.key === buttonKey);
        let animationNames = allAnimations.map(n => ({'text': n.name, 'value': n.name, 'params': n.params}));
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
                                      allSegments={allSegments}
                                      allAnimationNames={animationNames}
                                      onRemove={seg => this.props.onButtonChanged({
                                          ...button,
                                          segments: button.segments.filter(s => s.name !== seg.name || s.z !== seg.z)
                                      })}
                                      onSegmentChanged={seg => this.props.onButtonChanged({
                                          ...button,
                                          segments: button.segments.map(s => s.name === seg.name && s.z === seg.z ? seg : s)
                                      })}/>))}

                <div className='button-editor-ok-cancel-container'>
                    <LedSegmentChooser allSegments={allSegments}
                                       userSegments={userSegments}
                                       checkedSegmentNames={this.state.selectedSegments}
                                       onOk={() => {
                                           this.preSegChooserButton = undefined;
                                           this.setState({selectedSegments: []});
                                       }}
                                       onCancel={() => {
                                           if (this.preSegChooserButton !== undefined) {
                                               this.props.onButtonChanged(this.preSegChooserButton);
                                               this.preSegChooserButton = undefined;
                                           }
                                           this.setState({selectedSegments: []});
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
                        this.props.onButtonChanged(this.preEditButton);
                        this.props.history.goBack();
                    }} content='Cancel'/>
                </div>
            </Fragment>
        </div>
    }
}

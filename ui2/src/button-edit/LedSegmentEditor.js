import * as React from "react";
import {Dropdown, Image, Label, Segment} from "semantic-ui-react";
import './LedSegmentEditor.css';
import ColourEditor from "./ColourEditor";
import PropTypes from 'prop-types';

LedSegmentEditor.propTypes = {
    // the segment we are editing
    segment: PropTypes.shape({
        segment: PropTypes.string.isRequired,
        animation: PropTypes.string.isRequired,
        params: PropTypes.string
    }).isRequired,

    // list of all possible animation names
    allAnimationNames: PropTypes.arrayOf(PropTypes.shape({
        text: PropTypes.string,
        value: PropTypes.string
    })),

    // function to call, with new animation name, TODO generalize to send new segment?
    onAnimationNameChange: PropTypes.func.isRequired
};

export function LedSegmentEditor(props) {
    return <Segment attached color='blue' className='lse-container'>
        <Label content={props.segment.segment}
               onRemove={() => props.onRemove(props.segment)}
               color='blue'
               attached='top left'/>
        <div>
            <Image className='lse-image'
                   size='tiny'
                   inline
                   src={"/segment-icons/" + encodeURIComponent(props.segment.segment) + ".svg"}/>
        </div>
        <div>
            <Dropdown text={props.segment.animation}
                      options={props.allAnimationNames}
                      onChange={props.onAnimationNameChange}
                      className='lse-animation-name'/>
            <div className='lse-parameters'>
                <ColourEditor colour={'#' + props.segment.params}
                              onChangeColour={colour => {
                                  console.info('TODO Change Colour:' + colour)
                              }}/>
            </div>
        </div>
    </Segment>;
}

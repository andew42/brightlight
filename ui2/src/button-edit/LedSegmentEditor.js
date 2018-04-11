import * as React from "react";
import {Dropdown, Image, Label, Segment} from "semantic-ui-react";
import './LedSegmentEditor.css';

export function LedSegmentEditor(props) {
    return <Segment attached color='blue' className='led-segment-container'>
        <Label content={props.segment.segment}
               onRemove={() => props.onRemove(props.segment)}
               color='blue' attached='top left'/>
        <div>
            <Image className='led-segment-image'
                   size='tiny'
                   inline
                   src={"/segment-icons/" + encodeURIComponent(props.segment.segment) + ".svg"}/>
        </div>
        <div>
            <Dropdown text={props.segment.animation}
                      options={props.animationNames}
                      onChange={props.onAnimationNameChange}
                      className='led-segment-animation-name'/>
            <div className='led-segment-param'>{props.segment.params}</div>
        </div>
    </Segment>;
}

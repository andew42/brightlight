import * as React from "react";
import {Dropdown, Image, Label, Segment} from "semantic-ui-react";
import './LedSegmentEditor.css';
import ColourEditor from "./ColourEditor";
import Colour from "../colour/Colour";

export function LedSegmentEditor(props) {
    return <Segment attached color='blue' className='lse-container'>
        <Label content={props.segment.name}
               onRemove={() => props.onRemove(props.segment)}
               color='blue'
               attached='top left'/>
        <div>
            <Image className='lse-image'
                   size='tiny'
                   inline
                   src={"/segment-icons/" + encodeURIComponent(props.segment.name) + ".svg"}/>
        </div>
        <div>
            <Dropdown text={props.segment.animation}
                      options={props.allAnimationNames}
                      onChange={e => props.onSegmentChanged({
                          ...props.segment,
                          animation: e.target.textContent
                      })}
                      className='lse-animation-name'/>
            <div className='lse-parameters'>
                <ColourEditor colour={new Colour(props.segment.params)}
                              onColourChanged={colour => props.onSegmentChanged({
                                  ...props.segment,
                                  params: colour
                              })}/>
            </div>
        </div>
    </Segment>;
}

import * as React from "react";
import {Dropdown, Image, Label, Segment} from "semantic-ui-react";
import './LedSegmentEditor.css';
import ColourEditor from "./ColourEditor";
import Colour from "../colour/Colour";
import {RangeEditor} from "./RangeEditor";
import {CheckboxEditor} from "./CheckboxEditor";

export function LedSegmentEditor(props) {
    let segmentDefinition = props.allSegments.find(x => x.name === props.segment.name);
    let icon = <div style={{width:0}}></div>;
    if (segmentDefinition !== undefined && segmentDefinition.icon) {
        icon = <Image className='lse-image'
                          size='tiny'
                          inline
                          src={"/segment-icons/" + encodeURIComponent(props.segment.name) + ".svg"}/>
}
    return <Segment attached color='blue' className='lse-container'>
        <Label content={props.segment.name}
               onRemove={() => props.onRemove(props.segment)}
               color='blue'
               attached='top left'
               size='large'/>
        <div className='lse-animation-name'>
            <Dropdown value={props.segment.animation}
                      options={props.allAnimationNames}
                      onChange={(_, d) => {
                          let a = props.allAnimationNames.find(x => x.text === d.value);
                          props.onSegmentChanged({
                              ...props.segment,
                              animation: a.text,
                              params: a.params
                          })
                      }}/>
        </div>
        {icon}
        <div className='lse-parameters'>
            {props.segment.params.map(p => <LedSegmentParam
                key={p.key}
                param={p}
                onParamChanged={np => {
                    props.onSegmentChanged({
                        ...props.segment,
                        params: props.segment.params.map(p => p.key === np.key ? np : p)
                    })
                }}/>)}
        </div>
    </Segment>;
}

function LedSegmentParam(props) {
    switch (props.param.type) {
        case "colour":
            return <ColourEditor
                colour={new Colour(props.param.value)}
                label={props.param.label}
                onColourChanged={colour => props.onParamChanged({
                    ...props.param,
                    value: colour
                })}/>;
        case"range":
            return <RangeEditor
                label={props.param.label}
                min={props.param.min}
                max={props.param.max}
                value={props.param.value}
                onPosChanged={pos => props.onParamChanged({
                    ...props.param,
                    value: pos
                })}/>;
        case "checkbox":
            return <CheckboxEditor
                label={props.param.label}
                checked={props.param.value}
                onChange={newCheckState => props.onParamChanged({
                    ...props.param,
                    value: newCheckState
                })}/>;
        default:
            return <div>"Unknown Param: " + props.param.type</div>;
    }
}

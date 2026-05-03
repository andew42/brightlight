import * as React from "react";
import './LedSegmentEditor.css';
import ColourEditor from "./ColourEditor";
import Colour from "../colour/Colour";
import {RangeEditor} from "./RangeEditor";
import {CheckboxEditor} from "./CheckboxEditor";

export function LedSegmentEditor(props) {
    return (
        <fieldset className='lse-container'>
            <legend>
                <span>{props.segment.name}</span>
                <button className='lse-remove'
                        title='Remove'
                        onClick={() => props.onRemove(props.segment)}>×</button>
            </legend>
            <div className='lse-row'>
                <img className='lse-image'
                     alt=''
                     src={"/segment-icons/" + encodeURIComponent(props.segment.name) + ".svg"}
                     onError={e => e.target.style.display = 'none'}/>
                <div className='lse-animation-name'>
                    <select value={props.segment.animation}
                            onChange={e => {
                                let a = props.allAnimationNames.find(x => x.text === e.target.value);
                                props.onSegmentChanged({...props.segment, animation: a.text, params: a.params});
                            }}>
                        {props.allAnimationNames.map(a => <option key={a.value} value={a.value}>{a.text}</option>)}
                    </select>
                </div>
            </div>
            <div className='lse-parameters'>
                {props.segment.params.map(p => <LedSegmentParam
                    key={p.key}
                    param={p}
                    onParamChanged={np => {
                        props.onSegmentChanged({
                            ...props.segment,
                            params: props.segment.params.map(p => p.key === np.key ? np : p)
                        });
                    }}/>)}
            </div>
        </fieldset>
    );
}

function LedSegmentParam(props) {
    switch (props.param.type) {
        case "colour":
            return <ColourEditor
                colour={new Colour(props.param.value)}
                label={props.param.label}
                onColourChanged={colour => props.onParamChanged({...props.param, value: colour})}/>;
        case "range":
            return <RangeEditor
                label={props.param.label}
                min={props.param.min}
                max={props.param.max}
                value={props.param.value}
                onPosChanged={pos => props.onParamChanged({...props.param, value: pos})}/>;
        case "checkbox":
            return <CheckboxEditor
                label={props.param.label}
                checked={props.param.value}
                onChange={newCheckState => props.onParamChanged({...props.param, value: newCheckState})}/>;
        default:
            return <div>"Unknown Param: " + props.param.type</div>;
    }
}

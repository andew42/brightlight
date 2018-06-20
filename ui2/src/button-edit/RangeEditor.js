import * as React from "react";
import HorizontalSlider from "../buttons/HorizontalSlider";
import './RangeEditor.css';

export function RangeEditor(props) {
    return <HorizontalSlider min={props.min}
                             max={props.max}
                             pos={props.value}
                             onPosChange={p => props.onSpeedChanged(p)}
                             sliderColour='#2185d0'
                             className='range-editor'
                             label={props.label}/>
}

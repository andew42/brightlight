import * as React from "react";
import {Checkbox} from 'semantic-ui-react'

export function CheckboxEditor(props) {
    return <Checkbox
        label={props.label}
        checked={props.checked}
        onChange={(_, d) => props.onChange(d.checked)}/>
}

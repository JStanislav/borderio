import { useState } from "react";

interface InputFieldProps {
    title: string;

    initialValue: string;
    minLength?: number;
    maxLength?: number;

    submitText: string;
    onSubmit: Promise<((value: string) => void)>
}

export function InputField({title, initialValue, minLength, maxLength, onSubmit, submitText} : InputFieldProps) {
    const [value, setValue] = useState(initialValue);  
    
    const onClickButton = () => {
        if (minLength !== undefined && value.length < minLength) {
            console.error(`Input value must be at least ${minLength} characters long`);
            return;
        }
        if (maxLength !== undefined && value.length > maxLength) {
            console.error(`Input value must be at most ${maxLength} characters long`);
            return;
        }


        onSubmit.then((submitFunction) => {
            submitFunction(value);
        })
    }


    const onKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") {
            onClickButton();
        }
    }

    const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        e.preventDefault();
        setValue(e.target.value);   
    }

    return <div className="text-field-container">
        <span>{title}</span>
        <div className="text-field-and-button">
            <input type="text" value={value} onChange={onChange} onKeyDown={onKeyDown} />
            <button onClick={onClickButton}>{submitText}</button>
        </div>
    </div>
}
import { useState } from "react";
import "./inputfield.css";


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
    const [inputError, setInputError] = useState<string>(""); 
    
    // returns true if validation is ok and its respective error message if not
    const validateInput = (value: string): [boolean, string] => {
        if (minLength !== undefined && value.length < minLength) {
            return [false, `Input value must be at least ${minLength} characters long`];
        }
        if (maxLength !== undefined && value.length > maxLength) {
            return [false, `Input value must be at most ${maxLength} characters long`];
        }

        return [true, ""]
    }
    
    const onClickButton = () => {
        const [isValid, errorMessage] = validateInput(value);
        if (!isValid) {
            setInputError(errorMessage);
            return;
        }

        setInputError("");

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
        {inputError && <span className="input-error">{inputError}</span> }
        <div className="text-field-and-button">
            <input type="text" value={value} onChange={onChange} onKeyDown={onKeyDown} />
            <button onClick={onClickButton}>{submitText}</button>
        </div>
    </div>
}
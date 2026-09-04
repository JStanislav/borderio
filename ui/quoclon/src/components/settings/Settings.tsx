import { useAuth } from "../../contexts/auth-provider";
import { saveName } from "../../services/auth-service";
import { InputField } from "./InputField";


export const Settings = () => {
    const { user, setUser } = useAuth();

    const onSubmit = (value: string) => {
        if (user) {
            setUser({...user, name: value});
            saveName(value);
        }
    }

    return <div>
        <h1>Settings</h1>
        {user && <InputField 
                title="Name" 
                initialValue={user.name} 
                minLength={3} maxLength={13} 
                onSubmit={Promise.resolve(onSubmit)} 
                submitText="Save"
        />}
    </div>
}
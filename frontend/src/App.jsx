import './App.css';
import {Greet} from "../wailsjs/go/main/App";

function App() {
    // const [resultText, setResultText] = useState("Please enter your name below 👇");
    // const [name, setName] = useState('');
    // const updateName = (e) => setName(e.target.value);
    // const updateResultText = (result) => setResultText(result);

    function greet() {
        Greet(name).then(updateResultText);
    }

    return (
        <div id="App">
            <h1 className="text-3xl font-bold underline">
                Hello world!
            </h1>
        </div>
    )
}

export default App

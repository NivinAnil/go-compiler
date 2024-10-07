import axios from "axios"


async function executeCode(data) {
    const response = await axios.post("http://localhost:8080/api/v1/submission", data, {
        headers: {
          'Access-Control-Allow-Origin': '*',
        },
      });
    return response.data      
}

async function getExecution({ request_id }) {
  console.log(request_id);
  try {
    const response = await axios.get(`http://localhost:8081/api/v1/submissions/${request_id}`, {
      headers: {
        'Access-Control-Allow-Origin': '*', // This can be omitted if CORS is properly configured server-side
      },
    });
    console.log(response);
    return response.data;
  } catch (error) {
    console.error("Error fetching execution:", error);
  }
}


export { executeCode,getExecution }
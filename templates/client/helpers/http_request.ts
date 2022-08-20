import axios, { AxiosPromise } from 'axios'

type Method = 'get' | 'post' | 'put' | 'patch' | 'delete'

interface APIResponse {
  status: number
  data?: Record<string, unknown>
}

export const httpRequest = <T = APIResponse>(method: Method, url: string, body?: {}): AxiosPromise<T> => {
  const headers = {}
  const token = window.localStorage.getItem("token")
  if (token !== "") {
    headers["Authorization"] = `Bearer ${token}`
  }

  return new Promise((resolve, reject) => {
    axios({
      method,
      url,
      data: body,
      headers
    })
      .then((res) => resolve(res))
      .catch((err) => {
        window.console.error(err)
        if (err?.response?.status === 422) {
          window.location.reload()
        }
        reject(err)
      })
  })
}

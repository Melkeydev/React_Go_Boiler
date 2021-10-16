import React, { useState, useEffect, useCallback, useRef } from "react";
import { useForm } from "react-hook-form";
import axios from "axios";

export interface iLogin {
  username: string;
  password: string;
}

export const Login = ({ jwtProps }: { jwtProps: any }) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    getValues,
  } = useForm();
  const [loginState, setLoginState] = useState<iLogin>();

  const onSubmit = async (data: any) => {
    const { username, password } = data;

    const body = JSON.stringify({
      username,
      password,
    });

    const response = await axios.post("http://localhost:4000/v1/login", body);
    //TODO: Handle if nothing returns
    console.log(response.data.message);
    jwtProps(response.data.message);
    window.localStorage.setItem("jwt", JSON.stringify(response.data.response));
  };

  const handleChange = useCallback((e) => {
    const { id, value } = e.target;

    setLoginState((state: any) => ({
      ...state,
      [id]: value,
    }));
  }, []);

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <div className="min-h-screen bg-gray-100 flex">
        <div className="container mt-16 mx-auto max-w-md">
          <div className="py-12 p-10 bg-white rounded-xl">
            <div className="mb-2">
              <label className="mr-4 text-gray-700 font-bold inline-block mb-2">
                Username
              </label>

              <input
                type="text"
                className="border bg-gray-100 py-2 px-4 w-96 outline-none focus:ring-2 focus:ring-indigo-400 rounded"
                placeholder="Username"
                id="username"
                value={loginState?.username}
                {...register("username")}
                onChange={handleChange}
              />
            </div>

            <div className="mb-2">
              <label className="mr-4 text-gray-700 font-bold inline-block mb-2">
                Password
              </label>

              <input
                type="password"
                className="border bg-gray-100 py-2 px-4 w-96 outline-none focus:ring-2 focus:ring-indigo-400 rounded"
                placeholder="Password"
                id="password"
                value={loginState?.password}
                {...register("password", {
                  required: true,
                  maxLength: 50,
                  minLength: 8,
                })}
                onChange={handleChange}
              />
            </div>

            <button className="w-full mt-6 text-indigo-50 font-bold bg-indigo-600 py-3 rounded-md hover:bg-indigo-500 transition duration-300">
              Login
            </button>
          </div>
        </div>
      </div>
    </form>
  );
};

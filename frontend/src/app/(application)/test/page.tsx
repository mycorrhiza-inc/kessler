"use client"

import { useState } from "react";
import axios from "axios";

export default function Page() {
  const [filters, setFilters] = useState([]);
  
  const getFilters = async () => {
    try {
      const response = await axios.get(`http://localhost:3301/fugu/filters`);
      console.log(response.data);
      setFilters(response.data); // Assuming the API returns an array of filters
    } catch (error) {
      console.error('Error fetching filters:', error);
    }
  }

  return (
    <>
      <h2>Testing filters</h2>
      <button onClick={getFilters}>Get Filters</button>
      
      {filters.length > 0 && (
        <ul>
          {filters.map((filter, index) => (
            <li key={index}>{filter}</li>
          ))}
        </ul>
      )}
    </>
  );
}

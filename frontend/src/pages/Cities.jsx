import { useState, useEffect } from 'react'
import { getCities, getCountries } from '../api/resources'
import PageHeader from '../components/PageHeader'
import DataTable from '../components/DataTable'

export default function Cities() {
  const [cities, setCities] = useState([])
  const [countries, setCountries] = useState([])
  const [selectedCountryId, setSelectedCountryId] = useState('')
  const [searchTerm, setSearchTerm] = useState('') 
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    setLoading(true)
    Promise.all([getCities(), getCountries()])
      .then(([citiesData, countriesData]) => {
        const sortedCities = (citiesData || []).sort((a, b) => 
          a.name.localeCompare(b.name)
        );

        const sortedCountries = (countriesData || []).sort((a, b) => 
          a.name.localeCompare(b.name)
        );

        setCities(sortedCities);
        setCountries(sortedCountries);
      })
      .catch(err => {
        console.error("Error loading cities data:", err);
        setError("Failed to load cities or countries.");
      })
      .finally(() => setLoading(false));
  }, []);

  const filteredCities = cities.filter(city => {
    const matchesCountry = selectedCountryId ? city.country_id === parseInt(selectedCountryId) : true
    const matchesSearch = city.name.toLowerCase().includes(searchTerm.toLowerCase())
    return matchesCountry && matchesSearch
  })

  const columns = [
    { key: 'name', label: 'City Name' },
    { 
      key: 'country_id', 
      label: 'Country',
      render: (id) => countries.find(c => c.country_id === id)?.name || `ID: ${id}`
    },
    { 
      key: 'description', 
      label: 'Description',
      render: (text) => text || 'No description' 
    },
  ]

  if (error) return <div className="p-8 text-red-600">{error}</div>

  return (
    <div>
      <PageHeader 
        icon="🏙️" 
        title="Cities" 
        subtitle="Manage available destinations and their details" 
      />

      <div className="mb-6 flex flex-wrap gap-4 items-end">
        <div className="w-full max-w-xs">
          <label className="label">Search City</label>
          <input 
            type="text"
            className="input"
            placeholder="City name..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>

        <div className="w-full max-w-xs">
          <label className="label">Country</label>
          <select 
            className="input"
            value={selectedCountryId}
            onChange={(e) => setSelectedCountryId(e.target.value)}
          >
            <option value="">All Countries</option>
            {countries.map(c => (
              <option key={c.country_id} value={c.country_id}>{c.name}</option>
            ))}
          </select>
        </div>
      </div>

      <div className="card">
        <DataTable 
          columns={columns} 
          data={filteredCities} 
          loading={loading} 
        />
      </div>
    </div>
  )
}
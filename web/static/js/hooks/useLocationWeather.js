/**
 * usage:
 * import { getLocationWeather } from './hooks/useLocationWeather.js';
 * 
 * getLocationWeather().then(data => {
 *   // data = { city, temperature, icon, ... }
 * });
 */

const CACHE_KEY = 'locationCache';
const CACHE_DURATION = 24 * 60 * 60 * 1000; // 24 hours

// Helper: Map WMO weather codes to icons/emojis
function getWeatherIcon(code, isDay) {
    const icons = {
        0: isDay ? "ph-sun" : "ph-moon", // Clear sky
        1: isDay ? "ph-sun-dim" : "ph-moon", // Mainly clear
        2: "ph-cloud-sun", // Partly cloudy
        3: "ph-cloud", // Overcast
        45: "ph-fog", // Fog
        48: "ph-fog", // Deposits rime fog
        51: "ph-cloud-rain", // Drizzle: Light
        53: "ph-cloud-rain", // Drizzle: Moderate
        55: "ph-cloud-rain", // Drizzle: Dense
        61: "ph-cloud-rain", // Rain: Slight
        63: "ph-cloud-rain", // Rain: Moderate
        65: "ph-cloud-rain", // Rain: Heavy
        71: "ph-snowflake", // Snow: Slight
        73: "ph-snowflake", // Snow: Moderate
        75: "ph-snowflake", // Snow: Heavy
        80: "ph-cloud-rain", // Rain showers: Slight
        81: "ph-cloud-rain", // Rain showers: Moderate
        82: "ph-cloud-rain", // Rain showers: Violent
        95: "ph-lightning", // Thunderstorm: Slight or moderate
        96: "ph-lightning", // Thunderstorm with slight hail
        99: "ph-lightning", // Thunderstorm with heavy hail
    };
    return icons[code] || (isDay ? "ph-sun" : "ph-moon");
}

// 1. Check Local Cache
function checkCache() {
    try {
        const cached = localStorage.getItem(CACHE_KEY);
        if (cached) {
            const data = JSON.parse(cached);
            if (Date.now() - data.timestamp < CACHE_DURATION) {
                console.log('Weather loaded from cache');
                return data;
            }
        }
    } catch (e) {
        console.warn('Cache read error', e);
    }
    return null;
}

// 2. Fetch Weather from Open-Meteo
async function fetchWeatherFromAPI(lat, lon, city) {
    try {
        const url = `https://api.open-meteo.com/v1/forecast?latitude=${lat}&longitude=${lon}&current_weather=true`;
        const res = await fetch(url);
        if (!res.ok) throw new Error("Weather API Error");

        const data = await res.json();
        const weather = data.current_weather;

        const result = {
            city: city,
            latitude: lat,
            longitude: lon,
            temperature: Math.round(weather.temperature),
            icon: getWeatherIcon(weather.weathercode, weather.is_day),
            timestamp: Date.now()
        };

        // Save to cache
        localStorage.setItem(CACHE_KEY, JSON.stringify(result));
        return result;
    } catch (e) {
        console.error("Error fetching weather:", e);
        throw e;
    }
}

// 3. Fallback: IP-based Location
async function getIPLocation() {
    console.log("Attempting IP fallback...");
    const res = await fetch('https://ipapi.co/json/');
    if (!res.ok) throw new Error("IP API Error");
    const data = await res.json();

    // IP API returns 'city', 'latitude', 'longitude'
    return {
        lat: data.latitude,
        lon: data.longitude,
        city: data.city || "Unknown"
    };
}

// 4. Reverse Geocode (Lat/Lon -> City)
async function getCityName(lat, lon) {
    try {
        // Using geocode.maps.co (Free)
        const res = await fetch(`https://geocode.maps.co/reverse?lat=${lat}&lon=${lon}&api_key=`);
        // Note: Ideally needs an API key for higher limits, but works for low volume. 
        // Alternatives: nominatim.openstreetmap.org (requires User-Agent)

        // Let's use Nominatim as it is standard and doesn't require key for low usage if we set a proper header manually? 
        // Browsers set User-Agent automatically. 
        // We actully don't strictly need reverse geocoding if we are okay with just showing "Local" or if we use the IP fallback.
        // But the prompt asked us to use it.

        // The prompt specifically said: https://geocode.maps.co/reverse?lat={lat}&lon={lon}
        // Let's stick to that.

        const data = await res.json();
        return data.address?.city || data.address?.town || data.address?.village || "Location";
    } catch (e) {
        console.warn("Reverse geocode failed, using generic name");
        return "Local";
    }
}

// Main Hook Function
export async function getLocationWeather() {
    // A. Check Cache First
    const cached = checkCache();
    if (cached) return cached;

    return new Promise((resolve, reject) => {
        // B. Try Browser Geolocation
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                async (position) => {
                    try {
                        const { latitude, longitude } = position.coords;
                        // Fetch City Name
                        const city = await getCityName(latitude, longitude);
                        // Fetch Weather
                        const weatherData = await fetchWeatherFromAPI(latitude, longitude, city);
                        resolve(weatherData);
                    } catch (err) {
                        console.error("Geo success but fetch failed:", err);
                        // If logic fails inside success, try fallback
                        fallback();
                    }
                },
                (error) => {
                    console.warn("Geolocation denied/failed:", error);
                    fallback();
                },
                { enableHighAccuracy: true, timeout: 5000 }
            );
        } else {
            console.warn("Geolocation not supported");
            fallback();
        }

        // C. Fallback Routine
        async function fallback() {
            try {
                const { lat, lon, city } = await getIPLocation();
                const weatherData = await fetchWeatherFromAPI(lat, lon, city);
                resolve(weatherData);
            } catch (err) {
                console.error("All location methods failed:", err);
                reject({ error: "Location Unavailable" });
            }
        }
    });
}

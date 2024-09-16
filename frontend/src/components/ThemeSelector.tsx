const ThemeSelector = () => {
  const themes = [
    "light",
    "dark",
    "cupcake",
    "bumblebee",
    "emerald",
    "corporate",
    "synthwave",
    "retro",
    "cyberpunk",
    "valentine",
    "halloween",
    "garden",
    "forest",
    "aqua",
    "lofi",
    "pastel",
    "fantasy",
    "wireframe",
    "black",
    "luxury",
    "dracula",
    "cmyk",
    "autumn",
    "business",
    "acid",
    "lemonade",
    "night",
    "coffee",
    "winter",
    "dim",
    "nord",
    "sunset",
  ];
  return (
    <div className="join join-vertical">
      {themes.map((theme) => (
        <input
          type="radio"
          name="theme-buttons"
          className="btn theme-controller join-item"
          aria-label={theme.charAt(0).toUpperCase() + theme.slice(1)}
          value={theme.toLowerCase()}
        />
      ))}
      <input
        type="radio"
        name="theme-buttons"
        className="btn theme-controller join-item"
        aria-label="Aqua"
        value="aqua"
      />
    </div>
  );
};

export default ThemeSelector;

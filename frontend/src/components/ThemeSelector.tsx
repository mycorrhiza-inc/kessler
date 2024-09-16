const ThemeSelector = () => {
  const themes = [
    "light",
    "dark",
    "forest",
    "lemonade",
    "retro",
    // "colorblind",
    "cyberpunk",
    "valentine",
    "aqua",
  ];
  return (
    <div className="dropdown dropdown-hover">
      <div tabIndex={0} role="button" className="btn m-1">
        Theme
      </div>
      <ul className="menu dropdown-content bg-base-100 rounded-box z-[1] w-52 p-2 shadow">
        {themes.map((theme) => (
          <li key={theme}>
            <input
              type="radio"
              name="theme-dropdown"
              className="theme-controller btn btn-sm btn-block btn-ghost justify-start"
              aria-label={theme.charAt(0).toUpperCase() + theme.slice(1)}
              // Lowercase conversion probably not necessary
              value={theme.toLowerCase()}
            />
          </li>
        ))}
      </ul>
    </div>
  );
};

//<div className="dropdown mb-72" style={{ zIndex: 3001 }}>
//  <div tabIndex={0} role="button" className="btn m-1">
//    Theme
//  </div>
//  <ul
//    tabIndex={0}
//    className="dropdown-content bg-base-300 rounded-box z-[1] w-52 p-2 shadow-2xl"
//  >
//    {themes.map((theme) => (
//      <li key={theme}>
//        <input
//          type="radio"
//          name="theme-dropdown"
//          className="theme-controller btn btn-sm btn-block btn-ghost justify-start"
//          aria-label={theme.charAt(0).toUpperCase() + theme.slice(1)}
//          value={theme.toLowerCase()}
//        />
//      </li>
//    ))}
//    );
//  </ul>
//</div>
export default ThemeSelector;

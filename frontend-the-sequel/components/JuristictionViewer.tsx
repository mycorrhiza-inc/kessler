import React from "react";
import DeckGL from "@deck.gl/react";
import { MapViewState } from "@deck.gl/core";
import { LineLayer } from "@deck.gl/layers";

const INITIAL_VIEW_STATE: MapViewState = {
  longitude: -122.41669,
  latitude: 37.7853,
  zoom: 13,
};

type DataType = {
  from: [longitude: number, latitude: number];
  to: [longitude: number, latitude: number];
};

const TestViewer = () => {
  const layers = [
    new LineLayer<DataType>({
      id: "line-layer",
      data: "/path/to/data.json",
      getSourcePosition: (d: DataType) => d.from,
      getTargetPosition: (d: DataType) => d.to,
    }),
  ];

  return (
    <DeckGL initialViewState={INITIAL_VIEW_STATE} controller layers={layers} />
  );
};

export default TestViewer;

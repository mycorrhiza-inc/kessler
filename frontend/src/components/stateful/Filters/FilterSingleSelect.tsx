"use client"
import { TextPill } from '@/componenets/style/Pills/TextPills';
import React, { useState, useCallback, useMemo, useEffect } from 'react';

// Types for the component
interface FilterFieldDefinition {
  id: string;
  displayName: string;
  description: string;
  placeholder?: string;
  options?: Array<{
    value: string;
    label: string;
    disabled?: boolean;
  }>;
}

interface DynamicSingleSelectProps {
  fieldDefinition: FilterFieldDefinition;
  value: string;
  onChange: (value: string) => void;
  onFocus?: () => void;
  onBlur?: () => void;
  disabled?: boolean;
  className?: string;
  allowClear?: boolean;
  defaultValue?: string; // New prop for default value
  dynamicWidth?: boolean; // New prop to enable dynamic width
  minWidth?: string; // Minimum width when using dynamic width
  maxWidth?: string; // Maximum width when using dynamic width
}

/**
 * Enhanced Single Select Component with Search, Modern UI, Dynamic Width, and Default Value
 * Similar to DynamicMultiSelect but allows only one selection
 */
export function DynamicSingleSelect({
  fieldDefinition,
  value,
  onChange,
  onFocus,
  onBlur,
  disabled = false,
  className = "",
  allowClear = true,
  defaultValue = "", // Default to empty string
  dynamicWidth = false, // Default to false for backward compatibility
  minWidth = "200px", // Default minimum width
  maxWidth = "150px", // Default maximum width
}: DynamicSingleSelectProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [componentWidth, setComponentWidth] = useState<string>('100%');

  // Initialize with default value if value is empty and defaultValue is provided
  useEffect(() => {
    if (!value && defaultValue && fieldDefinition.options?.some(opt => opt.value === defaultValue)) {
      onChange(defaultValue);
    }
  }, [value, defaultValue, onChange, fieldDefinition.options]);

  const handleSelectionChange = useCallback((optionValue: string) => {
    onChange(optionValue);
    setIsOpen(false); // Close dropdown after selection
  }, [onChange]);

  const clearSelection = useCallback((event: React.MouseEvent) => {
    event.stopPropagation();
    onChange("");
  }, [onChange]);

  // Filter options based on search term
  const filteredOptions = useMemo(() => {
    if (!searchTerm) return fieldDefinition.options || [];
    return (fieldDefinition.options || []).filter(option =>
      option.label.toLowerCase().includes(searchTerm.toLowerCase()) ||
      option.value.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [fieldDefinition.options, searchTerm]);

  // Get selected option for display
  const selectedOption = useMemo(() => {
    return fieldDefinition.options?.find(opt => opt.value === value);
  }, [value, fieldDefinition.options]);

  // Calculate dynamic width based on content
  const calculateDynamicWidth = useCallback(() => {
    if (!dynamicWidth) return '100%';

    let calculatedWidth = '200px'; // Base width

    if (selectedOption) {
      // Calculate width based on selected option label length
      const labelLength = selectedOption.label.length;
      const displayNameLength = fieldDefinition.displayName.length;
      const totalLength = labelLength + displayNameLength + 10; // Add padding

      // Rough calculation: 8px per character + base padding
      const estimatedWidth = Math.max(200, totalLength * 8 + 80);
      calculatedWidth = `${Math.min(parseInt(maxWidth), estimatedWidth)}px`;
    } else if (fieldDefinition.placeholder) {
      // Calculate based on placeholder length
      const placeholderLength = fieldDefinition.placeholder.length;
      const estimatedWidth = Math.max(200, placeholderLength * 8 + 40);
      calculatedWidth = `${Math.min(parseInt(maxWidth), estimatedWidth)}px`;
    }

    return calculatedWidth;
  }, [selectedOption, fieldDefinition, dynamicWidth, maxWidth]);

  // Update component width when selection changes
  useEffect(() => {
    if (dynamicWidth) {
      const newWidth = calculateDynamicWidth();
      setComponentWidth(newWidth);
    }
  }, [selectedOption, dynamicWidth, calculateDynamicWidth]);

  const handleDropdownToggle = useCallback(() => {
    if (!disabled) {
      setIsOpen(!isOpen);
      if (!isOpen) {
        onFocus?.();
      } else {
        onBlur?.();
      }
    }
  }, [disabled, isOpen, onFocus, onBlur]);

  const handleSearchChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  }, []);

  const clearSearch = useCallback(() => {
    setSearchTerm('');
  }, []);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as Element;
      if (isOpen && !target.closest('.singleselect-container')) {
        setIsOpen(false);
        onBlur?.();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isOpen, onBlur]);

  // Dynamic styles for width
  const containerStyle = dynamicWidth ? {
    width: componentWidth,
    minWidth: minWidth,
    maxWidth: maxWidth,
    transition: 'width 0.2s ease-in-out'
  } : {};

  return (
    <div
      className={`relative z-50 singleselect-container ${className} ${dynamicWidth ? '' : 'w-full'}`}
      style={containerStyle}
    >
      {/* Main Button */}
      <button
        type="button"
        className={`
          w-full min-h-[3rem] p-3 text-left border-2 rounded-lg transition-all duration-200 flex items-center justify-between
          ${isOpen ? 'border-primary-500 ring-2 ring-primary-200' : 'border-gray-300'}
          ${disabled ? 'bg-gray-100 cursor-not-allowed' : 'bg-base-100 hover:border-gray-400 cursor-pointer'}
          ${selectedOption ? 'bg-base-100' : ''}
        `}
        onClick={handleDropdownToggle}
        disabled={disabled}
        onKeyDown={(e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            handleDropdownToggle();
          }
        }}
      >
        <div className="flex-1 min-w-0">
          {!selectedOption ? (
            <span className="text-gray-500">{fieldDefinition.placeholder || 'Select an option...'}</span>
          ) : (
            <div className="flex items-center justify-center">
              <TextPill
                text={selectedOption.label}
                seed={selectedOption.value}
                removable
                onRemove={() => clearSelection}
                size="xs"
                className="m-0 mt-0 mb-0 flex-shrink-0 ml-2"
              />
              {/* <div className="font-medium text-sm text-gray-700"> */}
              {/*   {fieldDefinition.displayName}: */}
              {/* </div> */}
              {/* <span className="inline-flex items-center gap-1 px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm"> */}
              {/*   {selectedOption.label} */}
              {/*   {allowClear && ( */}
              {/*     <button */}
              {/*       type="button" */}
              {/*       onClick={clearSelection} */}
              {/*       className="hover:bg-blue-200 rounded-full p-0.5 transition-colors ml-1" */}
              {/*       aria-label={`Clear ${selectedOption.label}`} */}
              {/*     > */}
              {/*       <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"> */}
              {/*         <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /> */}
              {/*       </svg> */}
              {/*     </button> */}
              {/*   )} */}
              {/* </span> */}
            </div>
          )}
        </div>

        {/* Dropdown Arrow */}
        <svg
          className={`w-5 h-5 text-gray-400 transition-transform flex-shrink-0 ml-2 ${isOpen ? 'rotate-180' : ''}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {/* Dropdown Content */}
      {isOpen && (
        <div className="absolute top-full left-0 right-0 z-50 mt-1 bg-white border-2 border-gray-200 rounded-lg shadow-lg max-h-80 overflow-hidden">
          {/* Search Box */}
          <div className="p-3 border-b border-gray-200 bg-gray-50">
            <div className="relative">
              <input
                type="text"
                className="w-full pl-9 pr-9 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder={`Search ${fieldDefinition.displayName.toLowerCase()}...`}
                value={searchTerm}
                onChange={handleSearchChange}
                onClick={(e) => e.stopPropagation()}
              />
              <svg
                className="w-4 h-4 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              {searchTerm && (
                <button
                  type="button"
                  onClick={clearSearch}
                  className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              )}
            </div>
          </div>

          {/* Options List */}
          <div className="max-h-60 overflow-y-auto">
            {/* Clear Option (if allowClear and something is selected) */}
            {allowClear && selectedOption && (
              <div className="border-b border-gray-100">
                <button
                  type="button"
                  onClick={() => handleSelectionChange("")}
                  className="w-full text-left p-3 hover:bg-gray-100 transition-colors text-gray-600 italic"
                >
                  <div className="flex items-center gap-3">
                    <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    <span>Clear selection</span>
                  </div>
                </button>
              </div>
            )}

            {filteredOptions.length === 0 ? (
              <div className="p-4 text-center text-gray-500 text-sm">
                {searchTerm ? 'No options match your search' : 'No options available'}
              </div>
            ) : (
              <div className="p-1">
                {filteredOptions.map((option) => {
                  const isSelected = value === option.value;

                  return (
                    <button
                      key={option.value}
                      type="button"
                      onClick={() => !option.disabled && handleSelectionChange(option.value)}
                      disabled={option.disabled}
                      className={`
                        w-full text-left p-3 rounded-md transition-colors flex items-center justify-between
                        ${option.disabled ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-100 cursor-pointer'}
                        ${isSelected ? 'bg-blue-50 border-l-4 border-blue-500' : ''}
                      `}
                    >
                      <div className="flex items-center gap-3 flex-1 min-w-0">
                        {/* Radio Button Indicator */}
                        <div className={`
                          w-4 h-4 rounded-full border-2 flex items-center justify-center flex-shrink-0
                          ${isSelected ? 'border-blue-500 bg-blue-500' : 'border-gray-300'}
                        `}>
                          {isSelected && (
                            <div className="w-2 h-2 bg-white rounded-full"></div>
                          )}
                        </div>
                        <span className="text-sm text-gray-900 truncate flex-1">{option.label}</span>
                      </div>

                      {/* Selected Indicator Pill */}
                      {isSelected && (
                        <span className="inline-flex items-center px-2 py-1 bg-blue-100 text-blue-800 rounded-full text-xs flex-shrink-0 ml-2">
                          Selected
                        </span>
                      )}
                    </button>
                  );
                })}
              </div>
            )}
          </div>

          {/* Footer */}
          <div className="p-3 border-t border-gray-200 bg-gray-50 text-xs text-gray-600">
            {selectedOption ? (
              <div className="flex justify-between items-center">
                <span>Selected: {selectedOption.label}</span>
                {allowClear && (
                  <button
                    type="button"
                    onClick={() => handleSelectionChange("")}
                    className="text-primary-600 hover:text-primary-800 underline font-medium"
                  >
                    Clear selection
                  </button>
                )}
              </div>
            ) : (
              <span>No selection made</span>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
